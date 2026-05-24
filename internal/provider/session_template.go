package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type SessionTemplatePart struct {
	Type     string `json:"type,omitempty"`
	Text     string `json:"text,omitempty"`
	Link     string `json:"link,omitempty"`
	ID       string `json:"id,omitempty"`
	FileName string `json:"filename,omitempty"`
}

type SessionTemplate struct {
	ID       string              `json:"id"`
	Name     string              `json:"name"`
	Language string              `json:"language,omitempty"`
	Category string              `json:"category,omitempty"`
	Header   SessionTemplatePart `json:"header,omitempty"`
	Body     string              `json:"body"`
	Footer   string              `json:"footer,omitempty"`
}

type TemplateLanguage struct {
	Code string `json:"code"`
}

type TemplateParameter struct {
	Type          string           `json:"type"`
	ParameterName string           `json:"parameter_name,omitempty"`
	Text          string           `json:"text,omitempty"`
	Image         *ImageContent    `json:"image,omitempty"`
	Document      *DocumentContent `json:"document,omitempty"`
	Action        *TemplateAction  `json:"action,omitempty"`
}

type TemplateAction struct {
	OrderDetails map[string]any `json:"order_details,omitempty"`
}

type TemplateComponent struct {
	Type       string              `json:"type"`
	SubType    string              `json:"sub_type,omitempty"`
	Index      any                 `json:"index,omitempty"`
	Parameters []TemplateParameter `json:"parameters,omitempty"`
}

type TemplateContent struct {
	Name       string              `json:"name"`
	Language   TemplateLanguage    `json:"language,omitempty"`
	Components []TemplateComponent `json:"components,omitempty"`
	Variables  map[string]string   `json:"variables,omitempty"`
}

type renderedTemplate struct {
	TemplateID     string
	HeaderType     string
	HeaderText     string
	Body           string
	Footer         string
	HeaderImage    *ImageContent
	HeaderDocument *DocumentContent
	OrderDetails   map[string]any
}

func (m *ClientManager) ListSessionTemplates(sessionID string) []SessionTemplate {
	sid := strings.TrimSpace(sessionID)
	m.templatesMu.RLock()
	items := cloneSessionTemplates(m.sessionTemplates[sid])
	m.templatesMu.RUnlock()
	sort.Slice(items, func(i, j int) bool {
		if items[i].Name == items[j].Name {
			return items[i].Language < items[j].Language
		}
		return items[i].Name < items[j].Name
	})
	return items
}

func (m *ClientManager) UpsertSessionTemplate(sessionID string, tpl SessionTemplate) (SessionTemplate, error) {
	sid := strings.TrimSpace(sessionID)
	if sid == "" {
		return SessionTemplate{}, errors.New("missing session id")
	}
	normalized, err := normalizeSessionTemplate(tpl)
	if err != nil {
		return SessionTemplate{}, err
	}

	m.templatesMu.Lock()
	prev := cloneSessionTemplates(m.sessionTemplates[sid])
	next := cloneSessionTemplates(prev)
	replaced := false
	for i := range next {
		if sameTemplateIdentity(next[i], normalized) {
			next[i] = normalized
			replaced = true
			break
		}
	}
	if !replaced {
		next = append(next, normalized)
	}
	m.sessionTemplates[sid] = next
	m.templatesMu.Unlock()

	if err := m.persistSessionTemplates(); err != nil {
		m.templatesMu.Lock()
		if len(prev) == 0 {
			delete(m.sessionTemplates, sid)
		} else {
			m.sessionTemplates[sid] = prev
		}
		m.templatesMu.Unlock()
		return SessionTemplate{}, err
	}
	return normalized, nil
}

func (m *ClientManager) DeleteSessionTemplate(sessionID, templateID string) (bool, error) {
	sid := strings.TrimSpace(sessionID)
	id := strings.TrimSpace(templateID)
	if sid == "" || id == "" {
		return false, errors.New("missing session id or template id")
	}

	m.templatesMu.Lock()
	prev := cloneSessionTemplates(m.sessionTemplates[sid])
	next := make([]SessionTemplate, 0, len(prev))
	removed := false
	for _, item := range prev {
		if item.ID == id {
			removed = true
			continue
		}
		next = append(next, item)
	}
	if removed {
		if len(next) == 0 {
			delete(m.sessionTemplates, sid)
		} else {
			m.sessionTemplates[sid] = next
		}
	}
	m.templatesMu.Unlock()

	if removed {
		if err := m.persistSessionTemplates(); err != nil {
			m.templatesMu.Lock()
			m.sessionTemplates[sid] = prev
			m.templatesMu.Unlock()
			return false, err
		}
	}
	return removed, nil
}

func (m *ClientManager) renderTemplateMessage(sessionID string, in *TemplateContent) (*renderedTemplate, error) {
	if in == nil {
		return nil, errors.New("template requires a template payload")
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, errors.New("template.name is required")
	}
	language := strings.TrimSpace(in.Language.Code)
	tpl, ok := m.findSessionTemplate(sessionID, name, language)
	if !ok {
		if language == "" {
			return nil, fmt.Errorf("template %q is not registered for session %q", name, sessionID)
		}
		return nil, fmt.Errorf("template %q language %q is not registered for session %q", name, language, sessionID)
	}

	componentVars := templateVariablesByComponent(in)
	headerVars := mergeTemplateVars(in.Variables, componentVars["header"])
	bodyVars := mergeTemplateVars(in.Variables, componentVars["body"])
	footerVars := mergeTemplateVars(in.Variables, componentVars["footer"])

	out := &renderedTemplate{
		TemplateID: tpl.ID,
		HeaderType: strings.ToLower(strings.TrimSpace(tpl.Header.Type)),
		HeaderText: renderTemplateText(tpl.Header.Text, headerVars),
		Body:       renderTemplateText(tpl.Body, bodyVars),
		Footer:     renderTemplateText(tpl.Footer, footerVars),
	}
	if out.HeaderType == "" && strings.TrimSpace(out.HeaderText) != "" {
		out.HeaderType = "text"
	}
	if orderDetails := templateOrderDetails(in); len(orderDetails) > 0 {
		out.OrderDetails = orderDetails
	}
	media := templateHeaderMedia(tpl, in)
	switch out.HeaderType {
	case "image":
		out.HeaderImage = media.image
	case "document":
		out.HeaderDocument = media.document
	case "", "text":
	default:
		return nil, fmt.Errorf("unsupported template header type: %s", out.HeaderType)
	}
	if strings.TrimSpace(out.Body) == "" {
		return nil, errors.New("rendered template body is empty")
	}
	return out, nil
}

func templateOrderDetails(in *TemplateContent) map[string]any {
	if in == nil {
		return nil
	}
	for _, component := range in.Components {
		if !strings.EqualFold(strings.TrimSpace(component.Type), "button") {
			continue
		}
		if subType := strings.TrimSpace(component.SubType); subType != "" && !strings.EqualFold(subType, "order_details") {
			continue
		}
		for _, param := range component.Parameters {
			if !strings.EqualFold(strings.TrimSpace(param.Type), "action") || param.Action == nil || len(param.Action.OrderDetails) == 0 {
				continue
			}
			return normalizeTemplateOrderDetails(param.Action.OrderDetails)
		}
	}
	return nil
}

func normalizeTemplateOrderDetails(in map[string]any) map[string]any {
	if len(in) == 0 {
		return nil
	}
	normalized := normalizeTemplateOrderValue(in)
	out, ok := normalized.(map[string]any)
	if !ok {
		return nil
	}
	return out
}

func normalizeTemplateOrderValue(value any) any {
	switch v := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))
		for key, child := range v {
			out[key] = normalizeTemplateOrderValueForKey(key, child)
		}
		return out
	case []any:
		out := make([]any, len(v))
		for i, child := range v {
			out[i] = normalizeTemplateOrderValue(child)
		}
		return out
	default:
		return v
	}
}

func normalizeTemplateOrderValueForKey(key string, value any) any {
	switch strings.ToLower(strings.TrimSpace(key)) {
	case "value", "offset", "quantity":
		if converted, ok := templateNumericValue(value); ok {
			return converted
		}
	}
	return normalizeTemplateOrderValue(value)
}

func templateNumericValue(value any) (int64, bool) {
	switch v := value.(type) {
	case json.Number:
		n, err := v.Int64()
		return n, err == nil
	case float64:
		return int64(v), true
	case float32:
		return int64(v), true
	case int:
		return int64(v), true
	case int64:
		return v, true
	case string:
		n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		return n, err == nil
	default:
		return 0, false
	}
}

func (m *ClientManager) findSessionTemplate(sessionID, name, language string) (SessionTemplate, bool) {
	sid := strings.TrimSpace(sessionID)
	name = strings.TrimSpace(name)
	language = strings.TrimSpace(language)
	m.templatesMu.RLock()
	items := cloneSessionTemplates(m.sessionTemplates[sid])
	m.templatesMu.RUnlock()
	for _, item := range items {
		if item.Name != name {
			continue
		}
		if language == "" || item.Language == "" || item.Language == language {
			return item, true
		}
	}
	return SessionTemplate{}, false
}

type templateHeaderMediaValue struct {
	image    *ImageContent
	document *DocumentContent
}

func templateHeaderMedia(tpl SessionTemplate, in *TemplateContent) templateHeaderMediaValue {
	var out templateHeaderMediaValue
	for _, component := range in.Components {
		if !strings.EqualFold(strings.TrimSpace(component.Type), "header") {
			continue
		}
		for _, param := range component.Parameters {
			switch strings.ToLower(strings.TrimSpace(param.Type)) {
			case "image":
				if param.Image != nil {
					img := *param.Image
					out.image = &img
				}
			case "document":
				if param.Document != nil {
					doc := *param.Document
					out.document = &doc
				}
			}
		}
	}
	if out.image == nil && strings.EqualFold(tpl.Header.Type, "image") && (tpl.Header.Link != "" || tpl.Header.ID != "") {
		out.image = &ImageContent{Link: tpl.Header.Link, ID: tpl.Header.ID}
	}
	if out.document == nil && strings.EqualFold(tpl.Header.Type, "document") && (tpl.Header.Link != "" || tpl.Header.ID != "") {
		out.document = &DocumentContent{Link: tpl.Header.Link, ID: tpl.Header.ID, FileName: tpl.Header.FileName}
	}
	return out
}

func templateVariablesByComponent(in *TemplateContent) map[string]map[string]string {
	out := map[string]map[string]string{}
	if in == nil {
		return out
	}
	for _, component := range in.Components {
		componentType := strings.ToLower(strings.TrimSpace(component.Type))
		if componentType == "" {
			continue
		}
		values := map[string]string{}
		position := 1
		for _, param := range component.Parameters {
			if strings.EqualFold(strings.TrimSpace(param.Type), "text") {
				text := strings.TrimSpace(param.Text)
				values[fmt.Sprintf("%d", position)] = text
				if name := strings.TrimSpace(param.ParameterName); name != "" {
					values[name] = text
				}
				position++
			}
		}
		if len(values) > 0 {
			out[componentType] = values
		}
	}
	return out
}

func mergeTemplateVars(global, component map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range global {
		key := strings.TrimSpace(k)
		if key != "" {
			out[key] = strings.TrimSpace(v)
		}
	}
	for k, v := range component {
		key := strings.TrimSpace(k)
		if key != "" {
			out[key] = strings.TrimSpace(v)
		}
	}
	return out
}

var templatePlaceholderRE = regexp.MustCompile(`\{\{\s*([A-Za-z0-9_.-]+)\s*\}\}`)

func renderTemplateText(text string, vars map[string]string) string {
	return strings.TrimSpace(templatePlaceholderRE.ReplaceAllStringFunc(text, func(match string) string {
		parts := templatePlaceholderRE.FindStringSubmatch(match)
		if len(parts) != 2 {
			return match
		}
		if value, ok := vars[parts[1]]; ok {
			return value
		}
		return match
	}))
}

func normalizeSessionTemplate(tpl SessionTemplate) (SessionTemplate, error) {
	tpl.ID = strings.TrimSpace(tpl.ID)
	tpl.Name = strings.TrimSpace(tpl.Name)
	tpl.Language = strings.TrimSpace(tpl.Language)
	tpl.Category = strings.TrimSpace(tpl.Category)
	tpl.Header.Type = strings.ToLower(strings.TrimSpace(tpl.Header.Type))
	tpl.Header.Text = strings.TrimSpace(tpl.Header.Text)
	tpl.Header.Link = strings.TrimSpace(tpl.Header.Link)
	tpl.Header.ID = strings.TrimSpace(tpl.Header.ID)
	tpl.Header.FileName = strings.TrimSpace(tpl.Header.FileName)
	tpl.Body = strings.TrimSpace(tpl.Body)
	tpl.Footer = strings.TrimSpace(tpl.Footer)
	if tpl.Name == "" {
		return SessionTemplate{}, errors.New("template name is required")
	}
	if tpl.Body == "" {
		return SessionTemplate{}, errors.New("template body is required")
	}
	switch tpl.Header.Type {
	case "", "text":
		if tpl.Header.Text != "" {
			tpl.Header.Type = "text"
		}
	case "image", "document":
	default:
		return SessionTemplate{}, fmt.Errorf("unsupported header type: %s", tpl.Header.Type)
	}
	if tpl.ID == "" {
		tpl.ID = tpl.Name
		if tpl.Language != "" {
			tpl.ID += "__" + tpl.Language
		}
	}
	return tpl, nil
}

func sameTemplateIdentity(left, right SessionTemplate) bool {
	return left.ID == right.ID || (left.Name == right.Name && left.Language == right.Language)
}

func cloneSessionTemplates(in []SessionTemplate) []SessionTemplate {
	if len(in) == 0 {
		return nil
	}
	out := make([]SessionTemplate, len(in))
	copy(out, in)
	return out
}

type sessionTemplateState struct {
	WorkerID string                       `json:"worker_id"`
	Sessions map[string][]SessionTemplate `json:"sessions"`
}

func loadSessionTemplateState(path, workerID string) (map[string][]SessionTemplate, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return map[string][]SessionTemplate{}, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string][]SessionTemplate{}, nil
		}
		return nil, err
	}
	st := sessionTemplateState{}
	if err := jsonUnmarshal(b, &st); err != nil {
		return nil, err
	}
	if strings.TrimSpace(st.WorkerID) != "" && strings.TrimSpace(workerID) != "" && strings.TrimSpace(st.WorkerID) != strings.TrimSpace(workerID) {
		return map[string][]SessionTemplate{}, nil
	}
	out := map[string][]SessionTemplate{}
	for sid, items := range st.Sessions {
		sid = normalizeSessionID(sid)
		if sid == "" {
			continue
		}
		for _, item := range items {
			tpl, err := normalizeSessionTemplate(item)
			if err == nil {
				out[sid] = append(out[sid], tpl)
			}
		}
	}
	return out, nil
}

func saveSessionTemplateState(path, workerID string, templates map[string][]SessionTemplate) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	payload := sessionTemplateState{
		WorkerID: strings.TrimSpace(workerID),
		Sessions: map[string][]SessionTemplate{},
	}
	for sid, items := range templates {
		sid = normalizeSessionID(sid)
		if sid == "" || len(items) == 0 {
			continue
		}
		payload.Sessions[sid] = cloneSessionTemplates(items)
	}
	b, err := jsonMarshalIndent(payload)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func (m *ClientManager) persistSessionTemplates() error {
	m.templatesMu.RLock()
	snapshot := make(map[string][]SessionTemplate, len(m.sessionTemplates))
	for sid, items := range m.sessionTemplates {
		snapshot[sid] = cloneSessionTemplates(items)
	}
	file := m.sessionTplFile
	m.templatesMu.RUnlock()
	return saveSessionTemplateState(file, m.workerID, snapshot)
}
