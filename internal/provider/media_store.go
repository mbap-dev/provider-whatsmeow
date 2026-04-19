package provider

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	ilog "your.org/provider-whatsmeow/internal/log"
)

type MediaStoreOptions struct {
	Mode          string
	LocalDir      string
	RegistryFile  string
	URLTTLSeconds int

	S3Bucket         string
	S3Region         string
	S3Endpoint       string
	S3AccessKeyID    string
	S3SecretAccess   string
	S3ForcePathStyle bool
	S3Prefix         string
}

type MediaSaveInput struct {
	MessagingProduct string
	MIMEType         string
	SHA256           string
	FileSize         int64
	Filename         string
}

type MediaMetadata struct {
	SessionID        string    `json:"-"`
	MediaID          string    `json:"-"`
	MessagingProduct string    `json:"messaging_product"`
	MIMEType         string    `json:"mime_type,omitempty"`
	SHA256           string    `json:"sha256,omitempty"`
	FileSize         int64     `json:"file_size,omitempty"`
	ID               string    `json:"id"`
	Filename         string    `json:"filename,omitempty"`
	RelativePath     string    `json:"-"`
	CreatedAt        time.Time `json:"-"`

	storageKey string
}

type mediaRecord struct {
	SessionID        string    `json:"session_id"`
	MediaID          string    `json:"media_id"`
	MessagingProduct string    `json:"messaging_product"`
	MIMEType         string    `json:"mime_type,omitempty"`
	SHA256           string    `json:"sha256,omitempty"`
	FileSize         int64     `json:"file_size,omitempty"`
	ID               string    `json:"id"`
	Filename         string    `json:"filename,omitempty"`
	RelativePath     string    `json:"relative_path"`
	StorageKey       string    `json:"storage_key"`
	CreatedAt        time.Time `json:"created_at"`
}

type mediaRegistryState struct {
	Records map[string]mediaRecord `json:"records"`
	Index   map[string]string      `json:"index"`
	Paths   map[string]string      `json:"paths"`
}

type mediaStore struct {
	mu sync.RWMutex

	mode         string
	localDir     string
	registryFile string
	urlTTL       time.Duration

	s3Client    *s3.Client
	s3Presigner *s3.PresignClient
	s3Bucket    string
	s3Prefix    string
	s3Endpoint  string
	s3PathStyle bool

	state mediaRegistryState
}

func newMediaStore(opts MediaStoreOptions) *mediaStore {
	store := &mediaStore{
		mode:         "local",
		localDir:     strings.TrimSpace(opts.LocalDir),
		registryFile: strings.TrimSpace(opts.RegistryFile),
		urlTTL:       time.Duration(defaultPositive(opts.URLTTLSeconds, 60*60*24*3)) * time.Second,
		state: mediaRegistryState{
			Records: map[string]mediaRecord{},
			Index:   map[string]string{},
			Paths:   map[string]string{},
		},
	}

	if store.localDir == "" {
		store.localDir = "./state/whatsmeow/_media"
	}
	if store.registryFile == "" {
		store.registryFile = filepath.Join(filepath.Dir(store.localDir), "_media_registry.json")
	}

	mode := strings.ToLower(strings.TrimSpace(opts.Mode))
	if mode == "" {
		mode = "auto"
	}
	if mode == "s3" || (mode == "auto" && hasS3Config(opts)) {
		if client, presigner, bucket, prefix, ok := buildS3Client(opts); ok {
			store.mode = "s3"
			store.s3Client = client
			store.s3Presigner = presigner
			store.s3Bucket = bucket
			store.s3Prefix = prefix
			store.s3Endpoint = strings.TrimSpace(opts.S3Endpoint)
			store.s3PathStyle = opts.S3ForcePathStyle
		}
	}

	store.logInit(mode, opts)
	if err := store.loadRegistry(); err != nil {
		ilog.Errorf("evt=media_registry load_failed file=%s err=%v", store.registryFile, err)
	} else {
		records, index, paths := store.registryCounts()
		ilog.Infof("evt=media_registry load_ok file=%s records=%d index=%d paths=%d", store.registryFile, records, index, paths)
	}
	return store
}

func hasS3Config(opts MediaStoreOptions) bool {
	return strings.TrimSpace(opts.S3Bucket) != "" && strings.TrimSpace(opts.S3AccessKeyID) != "" && strings.TrimSpace(opts.S3SecretAccess) != ""
}

func buildS3Client(opts MediaStoreOptions) (*s3.Client, *s3.PresignClient, string, string, bool) {
	bucket := strings.TrimSpace(opts.S3Bucket)
	accessKeyID := strings.TrimSpace(opts.S3AccessKeyID)
	secret := strings.TrimSpace(opts.S3SecretAccess)
	if bucket == "" || accessKeyID == "" || secret == "" {
		return nil, nil, "", "", false
	}

	region := strings.TrimSpace(opts.S3Region)
	if region == "" {
		region = "us-east-1"
	}

	cfg := aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider(accessKeyID, secret, ""),
	}

	endpoint := strings.TrimSpace(opts.S3Endpoint)
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = opts.S3ForcePathStyle
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})
	return client, s3.NewPresignClient(client), bucket, strings.Trim(strings.TrimSpace(opts.S3Prefix), "/"), true
}

func (s *mediaStore) logInit(requestedMode string, opts MediaStoreOptions) {
	region := strings.TrimSpace(opts.S3Region)
	if region == "" {
		region = "us-east-1"
	}
	switch s.mode {
	case "s3":
		ilog.Infof(
			"evt=media_store init requested_mode=%s resolved_mode=s3 bucket=%s region=%s endpoint=%s prefix=%s force_path_style=%t local_dir=%s registry_file=%s ttl=%s",
			requestedMode,
			s.s3Bucket,
			region,
			strings.TrimSpace(opts.S3Endpoint),
			s.s3Prefix,
			opts.S3ForcePathStyle,
			s.localDir,
			s.registryFile,
			s.urlTTL,
		)
		if strings.TrimSpace(s.s3Endpoint) != "" && !s.s3PathStyle {
			ilog.Infof(
				"evt=media_store s3_addressing mode=virtual_host bucket_host=%s hint=%q",
				virtualHostBucketEndpoint(s.s3Bucket, s.s3Endpoint),
				"Set S3_FORCE_PATH_STYLE=true if this endpoint does not resolve bucket subdomains",
			)
		}
	default:
		reason := "local_default"
		if requestedMode == "auto" {
			reason = "s3_config_absent"
			if hasS3Config(opts) {
				reason = "s3_init_fallback"
			}
		} else if requestedMode == "s3" {
			reason = "missing_s3_config"
		} else if requestedMode != "local" {
			reason = "unsupported_mode"
		}
		missing := missingS3ConfigFields(opts)
		ilog.Infof(
			"evt=media_store init requested_mode=%s resolved_mode=local reason=%s missing_s3_fields=%v local_dir=%s registry_file=%s ttl=%s",
			requestedMode,
			reason,
			missing,
			s.localDir,
			s.registryFile,
			s.urlTTL,
		)
		if requestedMode == "s3" && len(missing) > 0 {
			ilog.Errorf("evt=media_store s3_disabled requested_mode=s3 missing_s3_fields=%v", missing)
		}
	}
}

func missingS3ConfigFields(opts MediaStoreOptions) []string {
	var missing []string
	if strings.TrimSpace(opts.S3Bucket) == "" {
		missing = append(missing, "S3_BUCKET")
	}
	if strings.TrimSpace(opts.S3AccessKeyID) == "" {
		missing = append(missing, "S3_ACCESS_KEY_ID")
	}
	if strings.TrimSpace(opts.S3SecretAccess) == "" {
		missing = append(missing, "S3_SECRET_ACCESS_KEY")
	}
	return missing
}

func (s *mediaStore) registryCounts() (int, int, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.state.Records), len(s.state.Index), len(s.state.Paths)
}

func (s *mediaStore) Mode() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.mode
}

func (s *mediaStore) URLTTL() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.urlTTL
}

func (s *mediaStore) Save(ctx context.Context, sessionID, mediaID string, data []byte, input MediaSaveInput) (MediaMetadata, error) {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return MediaMetadata{}, errors.New("missing session id")
	}
	mediaID = normalizeMediaID(mediaID)
	if mediaID == "" {
		mediaID = generateMediaID()
	}
	if len(data) == 0 {
		return MediaMetadata{}, errors.New("empty media payload")
	}

	mimeType := mimeTypeBase(strings.TrimSpace(input.MIMEType))
	if mimeType == "" {
		mimeType = mimeTypeBase(http.DetectContentType(data))
	}
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	filename := strings.TrimSpace(input.Filename)
	relativePath := buildMediaRelativePath(sessionID, mediaID, mimeType, filename)
	if relativePath == "" {
		return MediaMetadata{}, errors.New("invalid media path")
	}

	storageKey := relativePath
	if s.mode == "s3" {
		storageKey = s.s3ObjectKey(relativePath)
	}

	entry := ilog.WithSession(sessionID).WithMessageID(mediaID)
	entry.Debug(
		"evt=media_save start mode=%s bytes=%d mime=%s filename=%q relative_path=%s storage_key=%s",
		s.mode,
		len(data),
		mimeType,
		filename,
		relativePath,
		storageKey,
	)

	if s.mode == "s3" {
		if err := s.writeS3(ctx, storageKey, data, mimeType); err != nil {
			entry.Error(
				"evt=media_save failed step=write_s3 mode=%s bucket=%s storage_key=%s bytes=%d mime=%s err=%v",
				s.mode,
				s.s3Bucket,
				storageKey,
				len(data),
				mimeType,
				err,
			)
			return MediaMetadata{}, err
		}
	} else {
		if err := s.writeLocal(relativePath, data); err != nil {
			entry.Error(
				"evt=media_save failed step=write_local mode=%s relative_path=%s full_path=%s bytes=%d mime=%s err=%v",
				s.mode,
				relativePath,
				filepath.Join(s.localDir, filepath.FromSlash(relativePath)),
				len(data),
				mimeType,
				err,
			)
			return MediaMetadata{}, err
		}
	}

	sz := input.FileSize
	if sz <= 0 {
		sz = int64(len(data))
	}
	messagingProduct := strings.TrimSpace(input.MessagingProduct)
	if messagingProduct == "" {
		messagingProduct = "whatsapp"
	}

	rec := mediaRecord{
		SessionID:        sessionID,
		MediaID:          mediaID,
		MessagingProduct: messagingProduct,
		MIMEType:         mimeType,
		SHA256:           strings.TrimSpace(input.SHA256),
		FileSize:         sz,
		ID:               cloudMediaID(sessionID, mediaID),
		Filename:         filename,
		RelativePath:     relativePath,
		StorageKey:       storageKey,
		CreatedAt:        time.Now().UTC(),
	}

	key := compositeMediaKey(sessionID, mediaID)
	s.mu.Lock()
	s.state.Records[key] = rec
	s.state.Index[mediaID] = key
	s.state.Paths[relativePath] = key
	err := s.persistRegistryLocked()
	s.mu.Unlock()
	if err != nil {
		entry.Error("evt=media_save failed step=persist_registry file=%s err=%v", s.registryFile, err)
		return MediaMetadata{}, err
	}
	entry.Info(
		"evt=media_save success mode=%s media_id=%s public_id=%s bytes=%d mime=%s relative_path=%s storage_key=%s filename=%q",
		s.mode,
		mediaID,
		rec.ID,
		sz,
		mimeType,
		relativePath,
		storageKey,
		filename,
	)
	return rec.toMetadata(), nil
}

func (s *mediaStore) Find(sessionID, mediaRef string) (MediaMetadata, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.findRecordLocked(sessionID, mediaRef)
	if !ok {
		return MediaMetadata{}, false
	}
	return rec.toMetadata(), true
}

func (s *mediaStore) FindByID(mediaID string) (MediaMetadata, bool) {
	mediaID = strings.TrimSpace(mediaID)
	if mediaID == "" {
		return MediaMetadata{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	key, ok := s.state.Index[mediaID]
	if !ok {
		return MediaMetadata{}, false
	}
	rec, ok := s.state.Records[key]
	if !ok {
		return MediaMetadata{}, false
	}
	return rec.toMetadata(), true
}

func (s *mediaStore) FindByRelativePath(relative string) (MediaMetadata, bool) {
	relative = normalizeRelativePath(relative)
	if relative == "" {
		return MediaMetadata{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	key, ok := s.state.Paths[relative]
	if !ok {
		return MediaMetadata{}, false
	}
	rec, ok := s.state.Records[key]
	if !ok {
		return MediaMetadata{}, false
	}
	return rec.toMetadata(), true
}

func (s *mediaStore) Load(ctx context.Context, sessionID, mediaRef string) ([]byte, MediaMetadata, error) {
	s.mu.RLock()
	rec, ok := s.findRecordLocked(sessionID, mediaRef)
	s.mu.RUnlock()
	if !ok {
		ilog.WithSession(sessionID).WithMessageID(mediaRef).Debug("evt=media_load not_found ref=%s", mediaRef)
		return nil, MediaMetadata{}, os.ErrNotExist
	}
	data, err := s.readRecord(ctx, rec)
	if err != nil {
		ilog.WithSession(rec.SessionID).WithMessageID(rec.MediaID).Error(
			"evt=media_load failed mode=%s storage_key=%s relative_path=%s err=%v",
			s.mode,
			rec.StorageKey,
			rec.RelativePath,
			err,
		)
		return nil, MediaMetadata{}, err
	}
	ilog.WithSession(rec.SessionID).WithMessageID(rec.MediaID).Info(
		"evt=media_load success mode=%s bytes=%d storage_key=%s relative_path=%s",
		s.mode,
		len(data),
		rec.StorageKey,
		rec.RelativePath,
	)
	return data, rec.toMetadata(), nil
}

func (s *mediaStore) OpenByRelativePath(ctx context.Context, relative string) (io.ReadCloser, MediaMetadata, error) {
	relative = normalizeRelativePath(relative)
	if relative == "" {
		return nil, MediaMetadata{}, os.ErrNotExist
	}

	s.mu.RLock()
	meta := MediaMetadata{RelativePath: relative}
	if key, ok := s.state.Paths[relative]; ok {
		if rec, okRec := s.state.Records[key]; okRec {
			meta = rec.toMetadata()
		}
	}
	mode := s.mode
	s.mu.RUnlock()

	if mode == "s3" {
		storageKey := meta.storageKey
		if strings.TrimSpace(storageKey) == "" {
			storageKey = s.s3ObjectKey(relative)
		}
		obj, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.s3Bucket),
			Key:    aws.String(storageKey),
		})
		if err != nil {
			ilog.WithSession(meta.SessionID).WithMessageID(meta.MediaID).Error(
				"evt=media_open failed mode=s3 bucket=%s storage_key=%s relative_path=%s err=%v",
				s.s3Bucket,
				storageKey,
				relative,
				err,
			)
			return nil, MediaMetadata{}, err
		}
		ilog.WithSession(meta.SessionID).WithMessageID(meta.MediaID).Info(
			"evt=media_open success mode=s3 bucket=%s storage_key=%s relative_path=%s",
			s.s3Bucket,
			storageKey,
			relative,
		)
		return obj.Body, meta, nil
	}

	fullPath := filepath.Join(s.localDir, filepath.FromSlash(relative))
	f, err := os.Open(fullPath)
	if err != nil {
		ilog.WithSession(meta.SessionID).WithMessageID(meta.MediaID).Error(
			"evt=media_open failed mode=local full_path=%s relative_path=%s err=%v",
			fullPath,
			relative,
			err,
		)
		return nil, MediaMetadata{}, err
	}
	ilog.WithSession(meta.SessionID).WithMessageID(meta.MediaID).Info(
		"evt=media_open success mode=local full_path=%s relative_path=%s",
		fullPath,
		relative,
	)
	return f, meta, nil
}

func (s *mediaStore) DirectURL(ctx context.Context, meta MediaMetadata) (string, error) {
	s.mu.RLock()
	mode := s.mode
	presigner := s.s3Presigner
	bucket := s.s3Bucket
	ttl := s.urlTTL
	s.mu.RUnlock()

	if mode != "s3" || presigner == nil {
		ilog.WithSession(meta.SessionID).WithMessageID(meta.MediaID).Debug(
			"evt=media_direct_url skipped mode=%s reason=unsupported_mode relative_path=%s",
			mode,
			meta.RelativePath,
		)
		return "", errors.New("direct URL not supported for local mode")
	}
	if strings.TrimSpace(meta.storageKey) == "" {
		ilog.WithSession(meta.SessionID).WithMessageID(meta.MediaID).Error(
			"evt=media_direct_url failed reason=missing_storage_key relative_path=%s",
			meta.RelativePath,
		)
		return "", errors.New("missing storage key")
	}

	in := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(meta.storageKey),
	}
	if strings.TrimSpace(meta.MIMEType) != "" {
		in.ResponseContentType = aws.String(meta.MIMEType)
	}
	out, err := presigner.PresignGetObject(ctx, in, func(o *s3.PresignOptions) {
		o.Expires = ttl
	})
	if err != nil {
		ilog.WithSession(meta.SessionID).WithMessageID(meta.MediaID).Error(
			"evt=media_direct_url failed mode=s3 bucket=%s storage_key=%s ttl=%s err=%v",
			bucket,
			meta.storageKey,
			ttl,
			err,
		)
		return "", err
	}
	ilog.WithSession(meta.SessionID).WithMessageID(meta.MediaID).Info(
		"evt=media_direct_url success mode=s3 bucket=%s storage_key=%s ttl=%s",
		bucket,
		meta.storageKey,
		ttl,
	)
	return out.URL, nil
}

func (s *mediaStore) writeLocal(relative string, data []byte) error {
	fullPath := filepath.Join(s.localDir, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(fullPath, data, 0o644)
}

func (s *mediaStore) writeS3(ctx context.Context, key string, data []byte, contentType string) error {
	in := &s3.PutObjectInput{
		Bucket: aws.String(s.s3Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}
	if strings.TrimSpace(contentType) != "" {
		in.ContentType = aws.String(contentType)
	}
	_, err := s.s3Client.PutObject(ctx, in)
	if err != nil && strings.TrimSpace(s.s3Endpoint) != "" && !s.s3PathStyle && strings.Contains(strings.ToLower(err.Error()), "no such host") {
		ilog.Errorf(
			"evt=media_store s3_hint action=put_object bucket=%s endpoint=%s bucket_host=%s suggestion=%q",
			s.s3Bucket,
			s.s3Endpoint,
			virtualHostBucketEndpoint(s.s3Bucket, s.s3Endpoint),
			"Enable S3_FORCE_PATH_STYLE=true or configure DNS/TLS for bucket subdomains",
		)
	}
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "signaturedoesnotmatch") {
		ilog.Errorf(
			"evt=media_store s3_hint action=put_object bucket=%s endpoint=%s force_path_style=%t suggestion=%q",
			s.s3Bucket,
			s.s3Endpoint,
			s.s3PathStyle,
			"Check S3 access key/secret, region, and whether the endpoint is the MinIO S3 API endpoint. If using a reverse proxy, preserve Host/path/Authorization headers without rewriting.",
		)
	}
	return err
}

func (s *mediaStore) readRecord(ctx context.Context, rec mediaRecord) ([]byte, error) {
	if s.mode == "s3" {
		obj, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.s3Bucket),
			Key:    aws.String(rec.StorageKey),
		})
		if err != nil {
			return nil, err
		}
		defer obj.Body.Close()
		return io.ReadAll(obj.Body)
	}
	fullPath := filepath.Join(s.localDir, filepath.FromSlash(rec.RelativePath))
	return os.ReadFile(fullPath)
}

func (s *mediaStore) s3ObjectKey(relative string) string {
	relative = normalizeRelativePath(relative)
	if strings.TrimSpace(s.s3Prefix) == "" {
		return relative
	}
	return path.Join(s.s3Prefix, relative)
}

func (s *mediaStore) loadRegistry() error {
	if strings.TrimSpace(s.registryFile) == "" {
		return nil
	}
	b, err := os.ReadFile(s.registryFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	st := mediaRegistryState{}
	if err := jsonUnmarshal(b, &st); err != nil {
		return err
	}
	s.mu.Lock()
	if st.Records != nil {
		s.state.Records = st.Records
	}
	if st.Index != nil {
		s.state.Index = st.Index
	}
	if st.Paths != nil {
		s.state.Paths = st.Paths
	}
	s.mu.Unlock()
	return nil
}

func (s *mediaStore) persistRegistryLocked() error {
	if strings.TrimSpace(s.registryFile) == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(s.registryFile), 0o755); err != nil {
		return err
	}
	payload, err := jsonMarshalIndent(s.state)
	if err != nil {
		return err
	}
	tmp := s.registryFile + ".tmp"
	if err := os.WriteFile(tmp, payload, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.registryFile)
}

func (s *mediaStore) findRecordLocked(sessionID, mediaRef string) (mediaRecord, bool) {
	sessionID = normalizeSessionID(sessionID)
	mediaRef = strings.TrimSpace(mediaRef)
	if mediaRef == "" {
		return mediaRecord{}, false
	}

	refSession, refMediaID := parseMediaRef(mediaRef)
	if refMediaID != "" {
		if sessionID != "" {
			if rec, ok := s.state.Records[compositeMediaKey(sessionID, refMediaID)]; ok {
				return rec, true
			}
		}
		if refSession != "" {
			if rec, ok := s.state.Records[compositeMediaKey(refSession, refMediaID)]; ok {
				return rec, true
			}
		}
	}

	if key, ok := s.state.Index[mediaRef]; ok {
		if rec, okRec := s.state.Records[key]; okRec {
			return rec, true
		}
	}
	if refMediaID != "" {
		if key, ok := s.state.Index[refMediaID]; ok {
			if rec, okRec := s.state.Records[key]; okRec {
				return rec, true
			}
		}
	}

	if sessionID != "" {
		if rec, ok := s.state.Records[compositeMediaKey(sessionID, mediaRef)]; ok {
			return rec, true
		}
	}

	return mediaRecord{}, false
}

func (r mediaRecord) toMetadata() MediaMetadata {
	return MediaMetadata{
		SessionID:        r.SessionID,
		MediaID:          r.MediaID,
		MessagingProduct: r.MessagingProduct,
		MIMEType:         r.MIMEType,
		SHA256:           r.SHA256,
		FileSize:         r.FileSize,
		ID:               r.ID,
		Filename:         r.Filename,
		RelativePath:     r.RelativePath,
		CreatedAt:        r.CreatedAt,
		storageKey:       r.StorageKey,
	}
}

func compositeMediaKey(sessionID, mediaID string) string {
	sessionID = normalizeSessionID(sessionID)
	mediaID = normalizeMediaID(mediaID)
	if sessionID == "" || mediaID == "" {
		return ""
	}
	return sessionID + "/" + mediaID
}

func parseMediaRef(ref string) (string, string) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return "", ""
	}
	parts := strings.Split(ref, "/")
	if len(parts) < 2 {
		return "", normalizeMediaID(ref)
	}
	sessionID := normalizeSessionID(parts[0])
	mediaID := normalizeMediaID(parts[len(parts)-1])
	return sessionID, mediaID
}

func buildMediaRelativePath(sessionID, mediaID, mimeType, filename string) string {
	sessionID = normalizePathToken(sessionID)
	mediaFile := normalizePathToken(mediaID)
	if sessionID == "" || mediaFile == "" {
		return ""
	}
	ext := mediaExtension(mimeType, filename)
	if ext == "" {
		ext = ".bin"
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return path.Join(sessionID, mediaFile+ext)
}

func mediaExtension(mimeType, filename string) string {
	filename = strings.TrimSpace(filename)
	if filename != "" {
		ext := strings.TrimSpace(filepath.Ext(filename))
		if ext != "" && len(ext) <= 12 {
			return ext
		}
	}
	mimeType = mimeTypeBase(mimeType)
	if mimeType == "" {
		return ""
	}
	ext, _ := mime.ExtensionsByType(mimeType)
	if len(ext) > 0 {
		return strings.TrimSpace(ext[0])
	}
	return ""
}

func normalizeMediaID(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if strings.Contains(v, "/") {
		parts := strings.Split(v, "/")
		v = strings.TrimSpace(parts[len(parts)-1])
	}
	return v
}

func normalizePathToken(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	var b strings.Builder
	for _, r := range v {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_' || r == '.':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	out := strings.Trim(b.String(), "._")
	if out == "" {
		return "media"
	}
	return out
}

func normalizeRelativePath(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	v = strings.ReplaceAll(v, "\\", "/")
	clean := path.Clean("/" + strings.TrimPrefix(v, "/"))
	clean = strings.TrimPrefix(clean, "/")
	if clean == "" || clean == "." {
		return ""
	}
	if strings.HasPrefix(clean, "../") || clean == ".." || strings.Contains(clean, "/../") {
		return ""
	}
	return clean
}

func virtualHostBucketEndpoint(bucket, endpoint string) string {
	bucket = strings.TrimSpace(bucket)
	endpoint = strings.TrimSpace(endpoint)
	if bucket == "" || endpoint == "" {
		return ""
	}
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimRight(endpoint, "/")
	if endpoint == "" {
		return ""
	}
	return bucket + "." + endpoint
}

func mimeTypeBase(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	return strings.TrimSpace(strings.SplitN(v, ";", 2)[0])
}

func generateMediaID() string {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return strings.ToUpper(hex.EncodeToString(buf))
}

func defaultPositive(v, def int) int {
	if v > 0 {
		return v
	}
	return def
}

func (m *ClientManager) StoreUploadedMedia(ctx context.Context, sessionID, filename, mimeType string, data []byte) (MediaMetadata, error) {
	if m == nil || m.mediaStore == nil {
		return MediaMetadata{}, errors.New("media store is not configured")
	}
	return m.mediaStore.Save(ctx, sessionID, "", data, MediaSaveInput{
		MessagingProduct: "whatsapp",
		MIMEType:         mimeType,
		Filename:         filename,
	})
}

func (m *ClientManager) StoreInboundMedia(ctx context.Context, sessionID, mediaID string, data []byte, input MediaSaveInput) (MediaMetadata, error) {
	if m == nil || m.mediaStore == nil {
		return MediaMetadata{}, errors.New("media store is not configured")
	}
	return m.mediaStore.Save(ctx, sessionID, mediaID, data, input)
}

func (m *ClientManager) GetMediaMetadata(sessionID, mediaRef string) (MediaMetadata, bool) {
	if m == nil || m.mediaStore == nil {
		return MediaMetadata{}, false
	}
	return m.mediaStore.Find(sessionID, mediaRef)
}

func (m *ClientManager) GetMediaMetadataByID(mediaID string) (MediaMetadata, bool) {
	if m == nil || m.mediaStore == nil {
		return MediaMetadata{}, false
	}
	return m.mediaStore.FindByID(mediaID)
}

func (m *ClientManager) MediaMode() string {
	if m == nil || m.mediaStore == nil {
		return "local"
	}
	return m.mediaStore.Mode()
}

func (m *ClientManager) MediaURLTTL() time.Duration {
	if m == nil || m.mediaStore == nil {
		return time.Hour
	}
	return m.mediaStore.URLTTL()
}

func CompatibleMediaID(meta MediaMetadata) string {
	if id := strings.TrimSpace(meta.MediaID); id != "" {
		return id
	}
	if _, mediaID := parseMediaRef(meta.ID); mediaID != "" {
		return mediaID
	}
	return normalizeMediaID(meta.ID)
}

func (m *ClientManager) MediaURL(ctx context.Context, baseURL, version string, meta MediaMetadata) string {
	mode := "local"
	if m != nil && m.mediaStore != nil {
		mode = m.mediaStore.Mode()
	}
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	version = strings.Trim(strings.TrimSpace(version), "/")
	if version == "" {
		version = "v1"
	}
	sessionID := strings.TrimSpace(meta.SessionID)
	if sessionID == "" {
		sessionID = normalizeSessionID(sessionFromCloudMediaID(meta.ID))
	}
	filePart := strings.TrimPrefix(strings.TrimSpace(meta.RelativePath), strings.TrimSpace(sessionID)+"/")
	filePart = strings.TrimPrefix(filePart, "/")
	if filePart == "" {
		filePart = normalizePathToken(meta.MediaID)
	}
	url := fmt.Sprintf("%s/%s/download/%s/%s", baseURL, version, sessionID, filePart)
	ilog.WithSession(meta.SessionID).WithMessageID(meta.MediaID).Info(
		"evt=media_url resolved mode=%s transport=proxy public_id=%s relative_path=%s version=%s compatibility=provider_auth",
		mode,
		meta.ID,
		meta.RelativePath,
		version,
	)
	return url
}

func (m *ClientManager) LoadMediaContent(ctx context.Context, sessionID, mediaRef string) ([]byte, MediaMetadata, error) {
	if m == nil || m.mediaStore == nil {
		return nil, MediaMetadata{}, errors.New("media store is not configured")
	}
	return m.mediaStore.Load(ctx, sessionID, mediaRef)
}

func (m *ClientManager) OpenMediaForDownload(ctx context.Context, sessionID, filePath string) (io.ReadCloser, MediaMetadata, error) {
	if m == nil || m.mediaStore == nil {
		return nil, MediaMetadata{}, errors.New("media store is not configured")
	}
	sessionID = normalizeSessionID(sessionID)
	relative := normalizeRelativePath(path.Join(sessionID, strings.TrimSpace(filePath)))
	if relative == "" {
		return nil, MediaMetadata{}, os.ErrNotExist
	}
	return m.mediaStore.OpenByRelativePath(ctx, relative)
}

func sessionFromCloudMediaID(id string) string {
	id = strings.TrimSpace(id)
	if id == "" {
		return ""
	}
	parts := strings.Split(id, "/")
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(parts[0])
}
