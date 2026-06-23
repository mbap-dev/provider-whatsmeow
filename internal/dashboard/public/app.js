(function(){const _=document.createElement("link").relList;if(_&&_.supports&&_.supports("modulepreload"))return;for(const C of document.querySelectorAll('link[rel="modulepreload"]'))T(C);new MutationObserver(C=>{for(const x of C)if(x.type==="childList")for(const M of x.addedNodes)M.tagName==="LINK"&&M.rel==="modulepreload"&&T(M)}).observe(document,{childList:!0,subtree:!0});function m(C){const x={};return C.integrity&&(x.integrity=C.integrity),C.referrerPolicy&&(x.referrerPolicy=C.referrerPolicy),C.crossOrigin==="use-credentials"?x.credentials="include":C.crossOrigin==="anonymous"?x.credentials="omit":x.credentials="same-origin",x}function T(C){if(C.ep)return;C.ep=!0;const x=m(C);fetch(C.href,x)}})();function Kf(f){return f&&f.__esModule&&Object.prototype.hasOwnProperty.call(f,"default")?f.default:f}var Yi={exports:{}},q={};var Ou;function Yf(){if(Ou)return q;Ou=1;var f=Symbol.for("react.element"),_=Symbol.for("react.portal"),m=Symbol.for("react.fragment"),T=Symbol.for("react.strict_mode"),C=Symbol.for("react.profiler"),x=Symbol.for("react.provider"),M=Symbol.for("react.context"),G=Symbol.for("react.forward_ref"),E=Symbol.for("react.suspense"),se=Symbol.for("react.memo"),he=Symbol.for("react.lazy"),ne=Symbol.iterator;function Z(d){return d===null||typeof d!="object"?null:(d=ne&&d[ne]||d["@@iterator"],typeof d=="function"?d:null)}var ut={isMounted:function(){return!1},enqueueForceUpdate:function(){},enqueueReplaceState:function(){},enqueueSetState:function(){}},De=Object.assign,ue={};function re(d,y,D){this.props=d,this.context=y,this.refs=ue,this.updater=D||ut}re.prototype.isReactComponent={},re.prototype.setState=function(d,y){if(typeof d!="object"&&typeof d!="function"&&d!=null)throw Error("setState(...): takes an object of state variables to update or a function which returns an object of state variables.");this.updater.enqueueSetState(this,d,y,"setState")},re.prototype.forceUpdate=function(d){this.updater.enqueueForceUpdate(this,d,"forceUpdate")};function ht(){}ht.prototype=re.prototype;function tt(d,y,D){this.props=d,this.context=y,this.refs=ue,this.updater=D||ut}var nt=tt.prototype=new ht;nt.constructor=tt,De(nt,re.prototype),nt.isPureReactComponent=!0;var we=Array.isArray,rt=Object.prototype.hasOwnProperty,Oe={current:null},Fe={key:!0,ref:!0,__self:!0,__source:!0};function Ue(d,y,D){var Y,R={},X=null,J=null;if(y!=null)for(Y in y.ref!==void 0&&(J=y.ref),y.key!==void 0&&(X=""+y.key),y)rt.call(y,Y)&&!Fe.hasOwnProperty(Y)&&(R[Y]=y[Y]);var b=arguments.length-2;if(b===1)R.children=D;else if(1<b){for(var ae=Array(b),Me=0;Me<b;Me++)ae[Me]=arguments[Me+2];R.children=ae}if(d&&d.defaultProps)for(Y in b=d.defaultProps,b)R[Y]===void 0&&(R[Y]=b[Y]);return{$$typeof:f,type:d,key:X,ref:J,props:R,_owner:Oe.current}}function Ct(d,y){return{$$typeof:f,type:d.type,key:y,ref:d.ref,props:d.props,_owner:d._owner}}function gt(d){return typeof d=="object"&&d!==null&&d.$$typeof===f}function yt(d){var y={"=":"=0",":":"=2"};return"$"+d.replace(/[=:]/g,function(D){return y[D]})}var Te=/\/+/g;function ze(d,y){return typeof d=="object"&&d!==null&&d.key!=null?yt(""+d.key):y.toString(36)}function Ce(d,y,D,Y,R){var X=typeof d;(X==="undefined"||X==="boolean")&&(d=null);var J=!1;if(d===null)J=!0;else switch(X){case"string":case"number":J=!0;break;case"object":switch(d.$$typeof){case f:case _:J=!0}}if(J)return J=d,R=R(J),d=Y===""?"."+ze(J,0):Y,we(R)?(D="",d!=null&&(D=d.replace(Te,"$&/")+"/"),Ce(R,y,D,"",function(Me){return Me})):R!=null&&(gt(R)&&(R=Ct(R,D+(!R.key||J&&J.key===R.key?"":(""+R.key).replace(Te,"$&/")+"/")+d)),y.push(R)),1;if(J=0,Y=Y===""?".":Y+":",we(d))for(var b=0;b<d.length;b++){X=d[b];var ae=Y+ze(X,b);J+=Ce(X,y,D,ae,R)}else if(ae=Z(d),typeof ae=="function")for(d=ae.call(d),b=0;!(X=d.next()).done;)X=X.value,ae=Y+ze(X,b++),J+=Ce(X,y,D,ae,R);else if(X==="object")throw y=String(d),Error("Objects are not valid as a React child (found: "+(y==="[object Object]"?"object with keys {"+Object.keys(d).join(", ")+"}":y)+"). If you meant to render a collection of children, use an array instead.");return J}function Be(d,y,D){if(d==null)return d;var Y=[],R=0;return Ce(d,Y,"","",function(X){return y.call(D,X,R++)}),Y}function U(d){if(d._status===-1){var y=d._result;y=y(),y.then(function(D){(d._status===0||d._status===-1)&&(d._status=1,d._result=D)},function(D){(d._status===0||d._status===-1)&&(d._status=2,d._result=D)}),d._status===-1&&(d._status=0,d._result=y)}if(d._status===1)return d._result.default;throw d._result}var Q={current:null},P={transition:null},A={ReactCurrentDispatcher:Q,ReactCurrentBatchConfig:P,ReactCurrentOwner:Oe};function w(){throw Error("act(...) is not supported in production builds of React.")}return q.Children={map:Be,forEach:function(d,y,D){Be(d,function(){y.apply(this,arguments)},D)},count:function(d){var y=0;return Be(d,function(){y++}),y},toArray:function(d){return Be(d,function(y){return y})||[]},only:function(d){if(!gt(d))throw Error("React.Children.only expected to receive a single React element child.");return d}},q.Component=re,q.Fragment=m,q.Profiler=C,q.PureComponent=tt,q.StrictMode=T,q.Suspense=E,q.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED=A,q.act=w,q.cloneElement=function(d,y,D){if(d==null)throw Error("React.cloneElement(...): The argument must be a React element, but you passed "+d+".");var Y=De({},d.props),R=d.key,X=d.ref,J=d._owner;if(y!=null){if(y.ref!==void 0&&(X=y.ref,J=Oe.current),y.key!==void 0&&(R=""+y.key),d.type&&d.type.defaultProps)var b=d.type.defaultProps;for(ae in y)rt.call(y,ae)&&!Fe.hasOwnProperty(ae)&&(Y[ae]=y[ae]===void 0&&b!==void 0?b[ae]:y[ae])}var ae=arguments.length-2;if(ae===1)Y.children=D;else if(1<ae){b=Array(ae);for(var Me=0;Me<ae;Me++)b[Me]=arguments[Me+2];Y.children=b}return{$$typeof:f,type:d.type,key:R,ref:X,props:Y,_owner:J}},q.createContext=function(d){return d={$$typeof:M,_currentValue:d,_currentValue2:d,_threadCount:0,Provider:null,Consumer:null,_defaultValue:null,_globalName:null},d.Provider={$$typeof:x,_context:d},d.Consumer=d},q.createElement=Ue,q.createFactory=function(d){var y=Ue.bind(null,d);return y.type=d,y},q.createRef=function(){return{current:null}},q.forwardRef=function(d){return{$$typeof:G,render:d}},q.isValidElement=gt,q.lazy=function(d){return{$$typeof:he,_payload:{_status:-1,_result:d},_init:U}},q.memo=function(d,y){return{$$typeof:se,type:d,compare:y===void 0?null:y}},q.startTransition=function(d){var y=P.transition;P.transition={};try{d()}finally{P.transition=y}},q.unstable_act=w,q.useCallback=function(d,y){return Q.current.useCallback(d,y)},q.useContext=function(d){return Q.current.useContext(d)},q.useDebugValue=function(){},q.useDeferredValue=function(d){return Q.current.useDeferredValue(d)},q.useEffect=function(d,y){return Q.current.useEffect(d,y)},q.useId=function(){return Q.current.useId()},q.useImperativeHandle=function(d,y,D){return Q.current.useImperativeHandle(d,y,D)},q.useInsertionEffect=function(d,y){return Q.current.useInsertionEffect(d,y)},q.useLayoutEffect=function(d,y){return Q.current.useLayoutEffect(d,y)},q.useMemo=function(d,y){return Q.current.useMemo(d,y)},q.useReducer=function(d,y,D){return Q.current.useReducer(d,y,D)},q.useRef=function(d){return Q.current.useRef(d)},q.useState=function(d){return Q.current.useState(d)},q.useSyncExternalStore=function(d,y,D){return Q.current.useSyncExternalStore(d,y,D)},q.useTransition=function(){return Q.current.useTransition()},q.version="18.3.1",q}var Fu;function Yu(){return Fu||(Fu=1,Yi.exports=Yf()),Yi.exports}var le=Yu();const Gf=Kf(le);var ro={},Gi={exports:{}},et={},Xi={exports:{}},qi={};var Mu;function Xf(){return Mu||(Mu=1,(function(f){function _(P,A){var w=P.length;P.push(A);e:for(;0<w;){var d=w-1>>>1,y=P[d];if(0<C(y,A))P[d]=A,P[w]=y,w=d;else break e}}function m(P){return P.length===0?null:P[0]}function T(P){if(P.length===0)return null;var A=P[0],w=P.pop();if(w!==A){P[0]=w;e:for(var d=0,y=P.length,D=y>>>1;d<D;){var Y=2*(d+1)-1,R=P[Y],X=Y+1,J=P[X];if(0>C(R,w))X<y&&0>C(J,R)?(P[d]=J,P[X]=w,d=X):(P[d]=R,P[Y]=w,d=Y);else if(X<y&&0>C(J,w))P[d]=J,P[X]=w,d=X;else break e}}return A}function C(P,A){var w=P.sortIndex-A.sortIndex;return w!==0?w:P.id-A.id}if(typeof performance=="object"&&typeof performance.now=="function"){var x=performance;f.unstable_now=function(){return x.now()}}else{var M=Date,G=M.now();f.unstable_now=function(){return M.now()-G}}var E=[],se=[],he=1,ne=null,Z=3,ut=!1,De=!1,ue=!1,re=typeof setTimeout=="function"?setTimeout:null,ht=typeof clearTimeout=="function"?clearTimeout:null,tt=typeof setImmediate<"u"?setImmediate:null;typeof navigator<"u"&&navigator.scheduling!==void 0&&navigator.scheduling.isInputPending!==void 0&&navigator.scheduling.isInputPending.bind(navigator.scheduling);function nt(P){for(var A=m(se);A!==null;){if(A.callback===null)T(se);else if(A.startTime<=P)T(se),A.sortIndex=A.expirationTime,_(E,A);else break;A=m(se)}}function we(P){if(ue=!1,nt(P),!De)if(m(E)!==null)De=!0,U(rt);else{var A=m(se);A!==null&&Q(we,A.startTime-P)}}function rt(P,A){De=!1,ue&&(ue=!1,ht(Ue),Ue=-1),ut=!0;var w=Z;try{for(nt(A),ne=m(E);ne!==null&&(!(ne.expirationTime>A)||P&&!yt());){var d=ne.callback;if(typeof d=="function"){ne.callback=null,Z=ne.priorityLevel;var y=d(ne.expirationTime<=A);A=f.unstable_now(),typeof y=="function"?ne.callback=y:ne===m(E)&&T(E),nt(A)}else T(E);ne=m(E)}if(ne!==null)var D=!0;else{var Y=m(se);Y!==null&&Q(we,Y.startTime-A),D=!1}return D}finally{ne=null,Z=w,ut=!1}}var Oe=!1,Fe=null,Ue=-1,Ct=5,gt=-1;function yt(){return!(f.unstable_now()-gt<Ct)}function Te(){if(Fe!==null){var P=f.unstable_now();gt=P;var A=!0;try{A=Fe(!0,P)}finally{A?ze():(Oe=!1,Fe=null)}}else Oe=!1}var ze;if(typeof tt=="function")ze=function(){tt(Te)};else if(typeof MessageChannel<"u"){var Ce=new MessageChannel,Be=Ce.port2;Ce.port1.onmessage=Te,ze=function(){Be.postMessage(null)}}else ze=function(){re(Te,0)};function U(P){Fe=P,Oe||(Oe=!0,ze())}function Q(P,A){Ue=re(function(){P(f.unstable_now())},A)}f.unstable_IdlePriority=5,f.unstable_ImmediatePriority=1,f.unstable_LowPriority=4,f.unstable_NormalPriority=3,f.unstable_Profiling=null,f.unstable_UserBlockingPriority=2,f.unstable_cancelCallback=function(P){P.callback=null},f.unstable_continueExecution=function(){De||ut||(De=!0,U(rt))},f.unstable_forceFrameRate=function(P){0>P||125<P?console.error("forceFrameRate takes a positive int between 0 and 125, forcing frame rates higher than 125 fps is not supported"):Ct=0<P?Math.floor(1e3/P):5},f.unstable_getCurrentPriorityLevel=function(){return Z},f.unstable_getFirstCallbackNode=function(){return m(E)},f.unstable_next=function(P){switch(Z){case 1:case 2:case 3:var A=3;break;default:A=Z}var w=Z;Z=A;try{return P()}finally{Z=w}},f.unstable_pauseExecution=function(){},f.unstable_requestPaint=function(){},f.unstable_runWithPriority=function(P,A){switch(P){case 1:case 2:case 3:case 4:case 5:break;default:P=3}var w=Z;Z=P;try{return A()}finally{Z=w}},f.unstable_scheduleCallback=function(P,A,w){var d=f.unstable_now();switch(typeof w=="object"&&w!==null?(w=w.delay,w=typeof w=="number"&&0<w?d+w:d):w=d,P){case 1:var y=-1;break;case 2:y=250;break;case 5:y=1073741823;break;case 4:y=1e4;break;default:y=5e3}return y=w+y,P={id:he++,callback:A,priorityLevel:P,startTime:w,expirationTime:y,sortIndex:-1},w>d?(P.sortIndex=w,_(se,P),m(E)===null&&P===m(se)&&(ue?(ht(Ue),Ue=-1):ue=!0,Q(we,w-d))):(P.sortIndex=y,_(E,P),De||ut||(De=!0,U(rt))),P},f.unstable_shouldYield=yt,f.unstable_wrapCallback=function(P){var A=Z;return function(){var w=Z;Z=A;try{return P.apply(this,arguments)}finally{Z=w}}}})(qi)),qi}var ju;function qf(){return ju||(ju=1,Xi.exports=Xf()),Xi.exports}var Au;function Zf(){if(Au)return et;Au=1;var f=Yu(),_=qf();function m(e){for(var t="https://reactjs.org/docs/error-decoder.html?invariant="+e,n=1;n<arguments.length;n++)t+="&args[]="+encodeURIComponent(arguments[n]);return"Minified React error #"+e+"; visit "+t+" for the full message or use the non-minified dev environment for full errors and additional helpful warnings."}var T=new Set,C={};function x(e,t){M(e,t),M(e+"Capture",t)}function M(e,t){for(C[e]=t,e=0;e<t.length;e++)T.add(t[e])}var G=!(typeof window>"u"||typeof window.document>"u"||typeof window.document.createElement>"u"),E=Object.prototype.hasOwnProperty,se=/^[:A-Z_a-z\u00C0-\u00D6\u00D8-\u00F6\u00F8-\u02FF\u0370-\u037D\u037F-\u1FFF\u200C-\u200D\u2070-\u218F\u2C00-\u2FEF\u3001-\uD7FF\uF900-\uFDCF\uFDF0-\uFFFD][:A-Z_a-z\u00C0-\u00D6\u00D8-\u00F6\u00F8-\u02FF\u0370-\u037D\u037F-\u1FFF\u200C-\u200D\u2070-\u218F\u2C00-\u2FEF\u3001-\uD7FF\uF900-\uFDCF\uFDF0-\uFFFD\-.0-9\u00B7\u0300-\u036F\u203F-\u2040]*$/,he={},ne={};function Z(e){return E.call(ne,e)?!0:E.call(he,e)?!1:se.test(e)?ne[e]=!0:(he[e]=!0,!1)}function ut(e,t,n,r){if(n!==null&&n.type===0)return!1;switch(typeof t){case"function":case"symbol":return!0;case"boolean":return r?!1:n!==null?!n.acceptsBooleans:(e=e.toLowerCase().slice(0,5),e!=="data-"&&e!=="aria-");default:return!1}}function De(e,t,n,r){if(t===null||typeof t>"u"||ut(e,t,n,r))return!0;if(r)return!1;if(n!==null)switch(n.type){case 3:return!t;case 4:return t===!1;case 5:return isNaN(t);case 6:return isNaN(t)||1>t}return!1}function ue(e,t,n,r,l,o,i){this.acceptsBooleans=t===2||t===3||t===4,this.attributeName=r,this.attributeNamespace=l,this.mustUseProperty=n,this.propertyName=e,this.type=t,this.sanitizeURL=o,this.removeEmptyString=i}var re={};"children dangerouslySetInnerHTML defaultValue defaultChecked innerHTML suppressContentEditableWarning suppressHydrationWarning style".split(" ").forEach(function(e){re[e]=new ue(e,0,!1,e,null,!1,!1)}),[["acceptCharset","accept-charset"],["className","class"],["htmlFor","for"],["httpEquiv","http-equiv"]].forEach(function(e){var t=e[0];re[t]=new ue(t,1,!1,e[1],null,!1,!1)}),["contentEditable","draggable","spellCheck","value"].forEach(function(e){re[e]=new ue(e,2,!1,e.toLowerCase(),null,!1,!1)}),["autoReverse","externalResourcesRequired","focusable","preserveAlpha"].forEach(function(e){re[e]=new ue(e,2,!1,e,null,!1,!1)}),"allowFullScreen async autoFocus autoPlay controls default defer disabled disablePictureInPicture disableRemotePlayback formNoValidate hidden loop noModule noValidate open playsInline readOnly required reversed scoped seamless itemScope".split(" ").forEach(function(e){re[e]=new ue(e,3,!1,e.toLowerCase(),null,!1,!1)}),["checked","multiple","muted","selected"].forEach(function(e){re[e]=new ue(e,3,!0,e,null,!1,!1)}),["capture","download"].forEach(function(e){re[e]=new ue(e,4,!1,e,null,!1,!1)}),["cols","rows","size","span"].forEach(function(e){re[e]=new ue(e,6,!1,e,null,!1,!1)}),["rowSpan","start"].forEach(function(e){re[e]=new ue(e,5,!1,e.toLowerCase(),null,!1,!1)});var ht=/[\-:]([a-z])/g;function tt(e){return e[1].toUpperCase()}"accent-height alignment-baseline arabic-form baseline-shift cap-height clip-path clip-rule color-interpolation color-interpolation-filters color-profile color-rendering dominant-baseline enable-background fill-opacity fill-rule flood-color flood-opacity font-family font-size font-size-adjust font-stretch font-style font-variant font-weight glyph-name glyph-orientation-horizontal glyph-orientation-vertical horiz-adv-x horiz-origin-x image-rendering letter-spacing lighting-color marker-end marker-mid marker-start overline-position overline-thickness paint-order panose-1 pointer-events rendering-intent shape-rendering stop-color stop-opacity strikethrough-position strikethrough-thickness stroke-dasharray stroke-dashoffset stroke-linecap stroke-linejoin stroke-miterlimit stroke-opacity stroke-width text-anchor text-decoration text-rendering underline-position underline-thickness unicode-bidi unicode-range units-per-em v-alphabetic v-hanging v-ideographic v-mathematical vector-effect vert-adv-y vert-origin-x vert-origin-y word-spacing writing-mode xmlns:xlink x-height".split(" ").forEach(function(e){var t=e.replace(ht,tt);re[t]=new ue(t,1,!1,e,null,!1,!1)}),"xlink:actuate xlink:arcrole xlink:role xlink:show xlink:title xlink:type".split(" ").forEach(function(e){var t=e.replace(ht,tt);re[t]=new ue(t,1,!1,e,"http://www.w3.org/1999/xlink",!1,!1)}),["xml:base","xml:lang","xml:space"].forEach(function(e){var t=e.replace(ht,tt);re[t]=new ue(t,1,!1,e,"http://www.w3.org/XML/1998/namespace",!1,!1)}),["tabIndex","crossOrigin"].forEach(function(e){re[e]=new ue(e,1,!1,e.toLowerCase(),null,!1,!1)}),re.xlinkHref=new ue("xlinkHref",1,!1,"xlink:href","http://www.w3.org/1999/xlink",!0,!1),["src","href","action","formAction"].forEach(function(e){re[e]=new ue(e,1,!1,e.toLowerCase(),null,!0,!0)});function nt(e,t,n,r){var l=re.hasOwnProperty(t)?re[t]:null;(l!==null?l.type!==0:r||!(2<t.length)||t[0]!=="o"&&t[0]!=="O"||t[1]!=="n"&&t[1]!=="N")&&(De(t,n,l,r)&&(n=null),r||l===null?Z(t)&&(n===null?e.removeAttribute(t):e.setAttribute(t,""+n)):l.mustUseProperty?e[l.propertyName]=n===null?l.type===3?!1:"":n:(t=l.attributeName,r=l.attributeNamespace,n===null?e.removeAttribute(t):(l=l.type,n=l===3||l===4&&n===!0?"":""+n,r?e.setAttributeNS(r,t,n):e.setAttribute(t,n))))}var we=f.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED,rt=Symbol.for("react.element"),Oe=Symbol.for("react.portal"),Fe=Symbol.for("react.fragment"),Ue=Symbol.for("react.strict_mode"),Ct=Symbol.for("react.profiler"),gt=Symbol.for("react.provider"),yt=Symbol.for("react.context"),Te=Symbol.for("react.forward_ref"),ze=Symbol.for("react.suspense"),Ce=Symbol.for("react.suspense_list"),Be=Symbol.for("react.memo"),U=Symbol.for("react.lazy"),Q=Symbol.for("react.offscreen"),P=Symbol.iterator;function A(e){return e===null||typeof e!="object"?null:(e=P&&e[P]||e["@@iterator"],typeof e=="function"?e:null)}var w=Object.assign,d;function y(e){if(d===void 0)try{throw Error()}catch(n){var t=n.stack.trim().match(/\n( *(at )?)/);d=t&&t[1]||""}return`
`+d+e}var D=!1;function Y(e,t){if(!e||D)return"";D=!0;var n=Error.prepareStackTrace;Error.prepareStackTrace=void 0;try{if(t)if(t=function(){throw Error()},Object.defineProperty(t.prototype,"props",{set:function(){throw Error()}}),typeof Reflect=="object"&&Reflect.construct){try{Reflect.construct(t,[])}catch(g){var r=g}Reflect.construct(e,[],t)}else{try{t.call()}catch(g){r=g}e.call(t.prototype)}else{try{throw Error()}catch(g){r=g}e()}}catch(g){if(g&&r&&typeof g.stack=="string"){for(var l=g.stack.split(`
`),o=r.stack.split(`
`),i=l.length-1,s=o.length-1;1<=i&&0<=s&&l[i]!==o[s];)s--;for(;1<=i&&0<=s;i--,s--)if(l[i]!==o[s]){if(i!==1||s!==1)do if(i--,s--,0>s||l[i]!==o[s]){var c=`
`+l[i].replace(" at new "," at ");return e.displayName&&c.includes("<anonymous>")&&(c=c.replace("<anonymous>",e.displayName)),c}while(1<=i&&0<=s);break}}}finally{D=!1,Error.prepareStackTrace=n}return(e=e?e.displayName||e.name:"")?y(e):""}function R(e){switch(e.tag){case 5:return y(e.type);case 16:return y("Lazy");case 13:return y("Suspense");case 19:return y("SuspenseList");case 0:case 2:case 15:return e=Y(e.type,!1),e;case 11:return e=Y(e.type.render,!1),e;case 1:return e=Y(e.type,!0),e;default:return""}}function X(e){if(e==null)return null;if(typeof e=="function")return e.displayName||e.name||null;if(typeof e=="string")return e;switch(e){case Fe:return"Fragment";case Oe:return"Portal";case Ct:return"Profiler";case Ue:return"StrictMode";case ze:return"Suspense";case Ce:return"SuspenseList"}if(typeof e=="object")switch(e.$$typeof){case yt:return(e.displayName||"Context")+".Consumer";case gt:return(e._context.displayName||"Context")+".Provider";case Te:var t=e.render;return e=e.displayName,e||(e=t.displayName||t.name||"",e=e!==""?"ForwardRef("+e+")":"ForwardRef"),e;case Be:return t=e.displayName||null,t!==null?t:X(e.type)||"Memo";case U:t=e._payload,e=e._init;try{return X(e(t))}catch{}}return null}function J(e){var t=e.type;switch(e.tag){case 24:return"Cache";case 9:return(t.displayName||"Context")+".Consumer";case 10:return(t._context.displayName||"Context")+".Provider";case 18:return"DehydratedFragment";case 11:return e=t.render,e=e.displayName||e.name||"",t.displayName||(e!==""?"ForwardRef("+e+")":"ForwardRef");case 7:return"Fragment";case 5:return t;case 4:return"Portal";case 3:return"Root";case 6:return"Text";case 16:return X(t);case 8:return t===Ue?"StrictMode":"Mode";case 22:return"Offscreen";case 12:return"Profiler";case 21:return"Scope";case 13:return"Suspense";case 19:return"SuspenseList";case 25:return"TracingMarker";case 1:case 0:case 17:case 2:case 14:case 15:if(typeof t=="function")return t.displayName||t.name||null;if(typeof t=="string")return t}return null}function b(e){switch(typeof e){case"boolean":case"number":case"string":case"undefined":return e;case"object":return e;default:return""}}function ae(e){var t=e.type;return(e=e.nodeName)&&e.toLowerCase()==="input"&&(t==="checkbox"||t==="radio")}function Me(e){var t=ae(e)?"checked":"value",n=Object.getOwnPropertyDescriptor(e.constructor.prototype,t),r=""+e[t];if(!e.hasOwnProperty(t)&&typeof n<"u"&&typeof n.get=="function"&&typeof n.set=="function"){var l=n.get,o=n.set;return Object.defineProperty(e,t,{configurable:!0,get:function(){return l.call(this)},set:function(i){r=""+i,o.call(this,i)}}),Object.defineProperty(e,t,{enumerable:n.enumerable}),{getValue:function(){return r},setValue:function(i){r=""+i},stopTracking:function(){e._valueTracker=null,delete e[t]}}}}function oe(e){e._valueTracker||(e._valueTracker=Me(e))}function dn(e){if(!e)return!1;var t=e._valueTracker;if(!t)return!0;var n=t.getValue(),r="";return e&&(r=ae(e)?e.checked?"true":"false":e.value),e=r,e!==n?(t.setValue(e),!0):!1}function pn(e){if(e=e||(typeof document<"u"?document:void 0),typeof e>"u")return null;try{return e.activeElement||e.body}catch{return e.body}}function Wt(e,t){var n=t.checked;return w({},t,{defaultChecked:void 0,defaultValue:void 0,value:void 0,checked:n??e._wrapperState.initialChecked})}function Vr(e,t){var n=t.defaultValue==null?"":t.defaultValue,r=t.checked!=null?t.checked:t.defaultChecked;n=b(t.value!=null?t.value:n),e._wrapperState={initialChecked:r,initialValue:n,controlled:t.type==="checkbox"||t.type==="radio"?t.checked!=null:t.value!=null}}function or(e,t){t=t.checked,t!=null&&nt(e,"checked",t,!1)}function mn(e,t){or(e,t);var n=b(t.value),r=t.type;if(n!=null)r==="number"?(n===0&&e.value===""||e.value!=n)&&(e.value=""+n):e.value!==""+n&&(e.value=""+n);else if(r==="submit"||r==="reset"){e.removeAttribute("value");return}t.hasOwnProperty("value")?ir(e,t.type,n):t.hasOwnProperty("defaultValue")&&ir(e,t.type,b(t.defaultValue)),t.checked==null&&t.defaultChecked!=null&&(e.defaultChecked=!!t.defaultChecked)}function Rn(e,t,n){if(t.hasOwnProperty("value")||t.hasOwnProperty("defaultValue")){var r=t.type;if(!(r!=="submit"&&r!=="reset"||t.value!==void 0&&t.value!==null))return;t=""+e._wrapperState.initialValue,n||t===e.value||(e.value=t),e.defaultValue=t}n=e.name,n!==""&&(e.name=""),e.defaultChecked=!!e._wrapperState.initialChecked,n!==""&&(e.name=n)}function ir(e,t,n){(t!=="number"||pn(e.ownerDocument)!==e)&&(n==null?e.defaultValue=""+e._wrapperState.initialValue:e.defaultValue!==""+n&&(e.defaultValue=""+n))}var vn=Array.isArray;function Qt(e,t,n,r){if(e=e.options,t){t={};for(var l=0;l<n.length;l++)t["$"+n[l]]=!0;for(n=0;n<e.length;n++)l=t.hasOwnProperty("$"+e[n].value),e[n].selected!==l&&(e[n].selected=l),l&&r&&(e[n].defaultSelected=!0)}else{for(n=""+b(n),t=null,l=0;l<e.length;l++){if(e[l].value===n){e[l].selected=!0,r&&(e[l].defaultSelected=!0);return}t!==null||e[l].disabled||(t=e[l])}t!==null&&(t.selected=!0)}}function Lt(e,t){if(t.dangerouslySetInnerHTML!=null)throw Error(m(91));return w({},t,{value:void 0,defaultValue:void 0,children:""+e._wrapperState.initialValue})}function Hr(e,t){var n=t.value;if(n==null){if(n=t.children,t=t.defaultValue,n!=null){if(t!=null)throw Error(m(92));if(vn(n)){if(1<n.length)throw Error(m(93));n=n[0]}t=n}t==null&&(t=""),n=t}e._wrapperState={initialValue:b(n)}}function Kr(e,t){var n=b(t.value),r=b(t.defaultValue);n!=null&&(n=""+n,n!==e.value&&(e.value=n),t.defaultValue==null&&e.defaultValue!==n&&(e.defaultValue=n)),r!=null&&(e.defaultValue=""+r)}function Yr(e){var t=e.textContent;t===e._wrapperState.initialValue&&t!==""&&t!==null&&(e.value=t)}function Gr(e){switch(e){case"svg":return"http://www.w3.org/2000/svg";case"math":return"http://www.w3.org/1998/Math/MathML";default:return"http://www.w3.org/1999/xhtml"}}function sr(e,t){return e==null||e==="http://www.w3.org/1999/xhtml"?Gr(t):e==="http://www.w3.org/2000/svg"&&t==="foreignObject"?"http://www.w3.org/1999/xhtml":e}var Ln,Xr=(function(e){return typeof MSApp<"u"&&MSApp.execUnsafeLocalFunction?function(t,n,r,l){MSApp.execUnsafeLocalFunction(function(){return e(t,n,r,l)})}:e})(function(e,t){if(e.namespaceURI!=="http://www.w3.org/2000/svg"||"innerHTML"in e)e.innerHTML=t;else{for(Ln=Ln||document.createElement("div"),Ln.innerHTML="<svg>"+t.valueOf().toString()+"</svg>",t=Ln.firstChild;e.firstChild;)e.removeChild(e.firstChild);for(;t.firstChild;)e.appendChild(t.firstChild)}});function hn(e,t){if(t){var n=e.firstChild;if(n&&n===e.lastChild&&n.nodeType===3){n.nodeValue=t;return}}e.textContent=t}var gn={animationIterationCount:!0,aspectRatio:!0,borderImageOutset:!0,borderImageSlice:!0,borderImageWidth:!0,boxFlex:!0,boxFlexGroup:!0,boxOrdinalGroup:!0,columnCount:!0,columns:!0,flex:!0,flexGrow:!0,flexPositive:!0,flexShrink:!0,flexNegative:!0,flexOrder:!0,gridArea:!0,gridRow:!0,gridRowEnd:!0,gridRowSpan:!0,gridRowStart:!0,gridColumn:!0,gridColumnEnd:!0,gridColumnSpan:!0,gridColumnStart:!0,fontWeight:!0,lineClamp:!0,lineHeight:!0,opacity:!0,order:!0,orphans:!0,tabSize:!0,widows:!0,zIndex:!0,zoom:!0,fillOpacity:!0,floodOpacity:!0,stopOpacity:!0,strokeDasharray:!0,strokeDashoffset:!0,strokeMiterlimit:!0,strokeOpacity:!0,strokeWidth:!0},a=["Webkit","ms","Moz","O"];Object.keys(gn).forEach(function(e){a.forEach(function(t){t=t+e.charAt(0).toUpperCase()+e.substring(1),gn[t]=gn[e]})});function u(e,t,n){return t==null||typeof t=="boolean"||t===""?"":n||typeof t!="number"||t===0||gn.hasOwnProperty(e)&&gn[e]?(""+t).trim():t+"px"}function L(e,t){e=e.style;for(var n in t)if(t.hasOwnProperty(n)){var r=n.indexOf("--")===0,l=u(n,t[n],r);n==="float"&&(n="cssFloat"),r?e.setProperty(n,l):e[n]=l}}var K=w({menuitem:!0},{area:!0,base:!0,br:!0,col:!0,embed:!0,hr:!0,img:!0,input:!0,keygen:!0,link:!0,meta:!0,param:!0,source:!0,track:!0,wbr:!0});function me(e,t){if(t){if(K[e]&&(t.children!=null||t.dangerouslySetInnerHTML!=null))throw Error(m(137,e));if(t.dangerouslySetInnerHTML!=null){if(t.children!=null)throw Error(m(60));if(typeof t.dangerouslySetInnerHTML!="object"||!("__html"in t.dangerouslySetInnerHTML))throw Error(m(61))}if(t.style!=null&&typeof t.style!="object")throw Error(m(62))}}function Vt(e,t){if(e.indexOf("-")===-1)return typeof t.is=="string";switch(e){case"annotation-xml":case"color-profile":case"font-face":case"font-face-src":case"font-face-uri":case"font-face-format":case"font-face-name":case"missing-glyph":return!1;default:return!0}}var io=null;function so(e){return e=e.target||e.srcElement||window,e.correspondingUseElement&&(e=e.correspondingUseElement),e.nodeType===3?e.parentNode:e}var ao=null,In=null,Dn=null;function Ji(e){if(e=Tr(e)){if(typeof ao!="function")throw Error(m(280));var t=e.stateNode;t&&(t=yl(t),ao(e.stateNode,e.type,t))}}function bi(e){In?Dn?Dn.push(e):Dn=[e]:In=e}function es(){if(In){var e=In,t=Dn;if(Dn=In=null,Ji(e),t)for(e=0;e<t.length;e++)Ji(t[e])}}function ts(e,t){return e(t)}function ns(){}var uo=!1;function rs(e,t,n){if(uo)return e(t,n);uo=!0;try{return ts(e,t,n)}finally{uo=!1,(In!==null||Dn!==null)&&(ns(),es())}}function ar(e,t){var n=e.stateNode;if(n===null)return null;var r=yl(n);if(r===null)return null;n=r[t];e:switch(t){case"onClick":case"onClickCapture":case"onDoubleClick":case"onDoubleClickCapture":case"onMouseDown":case"onMouseDownCapture":case"onMouseMove":case"onMouseMoveCapture":case"onMouseUp":case"onMouseUpCapture":case"onMouseEnter":(r=!r.disabled)||(e=e.type,r=!(e==="button"||e==="input"||e==="select"||e==="textarea")),e=!r;break e;default:e=!1}if(e)return null;if(n&&typeof n!="function")throw Error(m(231,t,typeof n));return n}var co=!1;if(G)try{var ur={};Object.defineProperty(ur,"passive",{get:function(){co=!0}}),window.addEventListener("test",ur,ur),window.removeEventListener("test",ur,ur)}catch{co=!1}function Zu(e,t,n,r,l,o,i,s,c){var g=Array.prototype.slice.call(arguments,3);try{t.apply(n,g)}catch(S){this.onError(S)}}var cr=!1,qr=null,Zr=!1,fo=null,Ju={onError:function(e){cr=!0,qr=e}};function bu(e,t,n,r,l,o,i,s,c){cr=!1,qr=null,Zu.apply(Ju,arguments)}function ec(e,t,n,r,l,o,i,s,c){if(bu.apply(this,arguments),cr){if(cr){var g=qr;cr=!1,qr=null}else throw Error(m(198));Zr||(Zr=!0,fo=g)}}function yn(e){var t=e,n=e;if(e.alternate)for(;t.return;)t=t.return;else{e=t;do t=e,(t.flags&4098)!==0&&(n=t.return),e=t.return;while(e)}return t.tag===3?n:null}function ls(e){if(e.tag===13){var t=e.memoizedState;if(t===null&&(e=e.alternate,e!==null&&(t=e.memoizedState)),t!==null)return t.dehydrated}return null}function os(e){if(yn(e)!==e)throw Error(m(188))}function tc(e){var t=e.alternate;if(!t){if(t=yn(e),t===null)throw Error(m(188));return t!==e?null:e}for(var n=e,r=t;;){var l=n.return;if(l===null)break;var o=l.alternate;if(o===null){if(r=l.return,r!==null){n=r;continue}break}if(l.child===o.child){for(o=l.child;o;){if(o===n)return os(l),e;if(o===r)return os(l),t;o=o.sibling}throw Error(m(188))}if(n.return!==r.return)n=l,r=o;else{for(var i=!1,s=l.child;s;){if(s===n){i=!0,n=l,r=o;break}if(s===r){i=!0,r=l,n=o;break}s=s.sibling}if(!i){for(s=o.child;s;){if(s===n){i=!0,n=o,r=l;break}if(s===r){i=!0,r=o,n=l;break}s=s.sibling}if(!i)throw Error(m(189))}}if(n.alternate!==r)throw Error(m(190))}if(n.tag!==3)throw Error(m(188));return n.stateNode.current===n?e:t}function is(e){return e=tc(e),e!==null?ss(e):null}function ss(e){if(e.tag===5||e.tag===6)return e;for(e=e.child;e!==null;){var t=ss(e);if(t!==null)return t;e=e.sibling}return null}var as=_.unstable_scheduleCallback,us=_.unstable_cancelCallback,nc=_.unstable_shouldYield,rc=_.unstable_requestPaint,ke=_.unstable_now,lc=_.unstable_getCurrentPriorityLevel,po=_.unstable_ImmediatePriority,cs=_.unstable_UserBlockingPriority,Jr=_.unstable_NormalPriority,oc=_.unstable_LowPriority,fs=_.unstable_IdlePriority,br=null,xt=null;function ic(e){if(xt&&typeof xt.onCommitFiberRoot=="function")try{xt.onCommitFiberRoot(br,e,void 0,(e.current.flags&128)===128)}catch{}}var _t=Math.clz32?Math.clz32:uc,sc=Math.log,ac=Math.LN2;function uc(e){return e>>>=0,e===0?32:31-(sc(e)/ac|0)|0}var el=64,tl=4194304;function fr(e){switch(e&-e){case 1:return 1;case 2:return 2;case 4:return 4;case 8:return 8;case 16:return 16;case 32:return 32;case 64:case 128:case 256:case 512:case 1024:case 2048:case 4096:case 8192:case 16384:case 32768:case 65536:case 131072:case 262144:case 524288:case 1048576:case 2097152:return e&4194240;case 4194304:case 8388608:case 16777216:case 33554432:case 67108864:return e&130023424;case 134217728:return 134217728;case 268435456:return 268435456;case 536870912:return 536870912;case 1073741824:return 1073741824;default:return e}}function nl(e,t){var n=e.pendingLanes;if(n===0)return 0;var r=0,l=e.suspendedLanes,o=e.pingedLanes,i=n&268435455;if(i!==0){var s=i&~l;s!==0?r=fr(s):(o&=i,o!==0&&(r=fr(o)))}else i=n&~l,i!==0?r=fr(i):o!==0&&(r=fr(o));if(r===0)return 0;if(t!==0&&t!==r&&(t&l)===0&&(l=r&-r,o=t&-t,l>=o||l===16&&(o&4194240)!==0))return t;if((r&4)!==0&&(r|=n&16),t=e.entangledLanes,t!==0)for(e=e.entanglements,t&=r;0<t;)n=31-_t(t),l=1<<n,r|=e[n],t&=~l;return r}function cc(e,t){switch(e){case 1:case 2:case 4:return t+250;case 8:case 16:case 32:case 64:case 128:case 256:case 512:case 1024:case 2048:case 4096:case 8192:case 16384:case 32768:case 65536:case 131072:case 262144:case 524288:case 1048576:case 2097152:return t+5e3;case 4194304:case 8388608:case 16777216:case 33554432:case 67108864:return-1;case 134217728:case 268435456:case 536870912:case 1073741824:return-1;default:return-1}}function fc(e,t){for(var n=e.suspendedLanes,r=e.pingedLanes,l=e.expirationTimes,o=e.pendingLanes;0<o;){var i=31-_t(o),s=1<<i,c=l[i];c===-1?((s&n)===0||(s&r)!==0)&&(l[i]=cc(s,t)):c<=t&&(e.expiredLanes|=s),o&=~s}}function mo(e){return e=e.pendingLanes&-1073741825,e!==0?e:e&1073741824?1073741824:0}function ds(){var e=el;return el<<=1,(el&4194240)===0&&(el=64),e}function vo(e){for(var t=[],n=0;31>n;n++)t.push(e);return t}function dr(e,t,n){e.pendingLanes|=t,t!==536870912&&(e.suspendedLanes=0,e.pingedLanes=0),e=e.eventTimes,t=31-_t(t),e[t]=n}function dc(e,t){var n=e.pendingLanes&~t;e.pendingLanes=t,e.suspendedLanes=0,e.pingedLanes=0,e.expiredLanes&=t,e.mutableReadLanes&=t,e.entangledLanes&=t,t=e.entanglements;var r=e.eventTimes;for(e=e.expirationTimes;0<n;){var l=31-_t(n),o=1<<l;t[l]=0,r[l]=-1,e[l]=-1,n&=~o}}function ho(e,t){var n=e.entangledLanes|=t;for(e=e.entanglements;n;){var r=31-_t(n),l=1<<r;l&t|e[r]&t&&(e[r]|=t),n&=~l}}var ie=0;function ps(e){return e&=-e,1<e?4<e?(e&268435455)!==0?16:536870912:4:1}var ms,go,vs,hs,gs,yo=!1,rl=[],Ht=null,Kt=null,Yt=null,pr=new Map,mr=new Map,Gt=[],pc="mousedown mouseup touchcancel touchend touchstart auxclick dblclick pointercancel pointerdown pointerup dragend dragstart drop compositionend compositionstart keydown keypress keyup input textInput copy cut paste click change contextmenu reset submit".split(" ");function ys(e,t){switch(e){case"focusin":case"focusout":Ht=null;break;case"dragenter":case"dragleave":Kt=null;break;case"mouseover":case"mouseout":Yt=null;break;case"pointerover":case"pointerout":pr.delete(t.pointerId);break;case"gotpointercapture":case"lostpointercapture":mr.delete(t.pointerId)}}function vr(e,t,n,r,l,o){return e===null||e.nativeEvent!==o?(e={blockedOn:t,domEventName:n,eventSystemFlags:r,nativeEvent:o,targetContainers:[l]},t!==null&&(t=Tr(t),t!==null&&go(t)),e):(e.eventSystemFlags|=r,t=e.targetContainers,l!==null&&t.indexOf(l)===-1&&t.push(l),e)}function mc(e,t,n,r,l){switch(t){case"focusin":return Ht=vr(Ht,e,t,n,r,l),!0;case"dragenter":return Kt=vr(Kt,e,t,n,r,l),!0;case"mouseover":return Yt=vr(Yt,e,t,n,r,l),!0;case"pointerover":var o=l.pointerId;return pr.set(o,vr(pr.get(o)||null,e,t,n,r,l)),!0;case"gotpointercapture":return o=l.pointerId,mr.set(o,vr(mr.get(o)||null,e,t,n,r,l)),!0}return!1}function _s(e){var t=_n(e.target);if(t!==null){var n=yn(t);if(n!==null){if(t=n.tag,t===13){if(t=ls(n),t!==null){e.blockedOn=t,gs(e.priority,function(){vs(n)});return}}else if(t===3&&n.stateNode.current.memoizedState.isDehydrated){e.blockedOn=n.tag===3?n.stateNode.containerInfo:null;return}}}e.blockedOn=null}function ll(e){if(e.blockedOn!==null)return!1;for(var t=e.targetContainers;0<t.length;){var n=wo(e.domEventName,e.eventSystemFlags,t[0],e.nativeEvent);if(n===null){n=e.nativeEvent;var r=new n.constructor(n.type,n);io=r,n.target.dispatchEvent(r),io=null}else return t=Tr(n),t!==null&&go(t),e.blockedOn=n,!1;t.shift()}return!0}function ws(e,t,n){ll(e)&&n.delete(t)}function vc(){yo=!1,Ht!==null&&ll(Ht)&&(Ht=null),Kt!==null&&ll(Kt)&&(Kt=null),Yt!==null&&ll(Yt)&&(Yt=null),pr.forEach(ws),mr.forEach(ws)}function hr(e,t){e.blockedOn===t&&(e.blockedOn=null,yo||(yo=!0,_.unstable_scheduleCallback(_.unstable_NormalPriority,vc)))}function gr(e){function t(l){return hr(l,e)}if(0<rl.length){hr(rl[0],e);for(var n=1;n<rl.length;n++){var r=rl[n];r.blockedOn===e&&(r.blockedOn=null)}}for(Ht!==null&&hr(Ht,e),Kt!==null&&hr(Kt,e),Yt!==null&&hr(Yt,e),pr.forEach(t),mr.forEach(t),n=0;n<Gt.length;n++)r=Gt[n],r.blockedOn===e&&(r.blockedOn=null);for(;0<Gt.length&&(n=Gt[0],n.blockedOn===null);)_s(n),n.blockedOn===null&&Gt.shift()}var On=we.ReactCurrentBatchConfig,ol=!0;function hc(e,t,n,r){var l=ie,o=On.transition;On.transition=null;try{ie=1,_o(e,t,n,r)}finally{ie=l,On.transition=o}}function gc(e,t,n,r){var l=ie,o=On.transition;On.transition=null;try{ie=4,_o(e,t,n,r)}finally{ie=l,On.transition=o}}function _o(e,t,n,r){if(ol){var l=wo(e,t,n,r);if(l===null)Mo(e,t,r,il,n),ys(e,r);else if(mc(l,e,t,n,r))r.stopPropagation();else if(ys(e,r),t&4&&-1<pc.indexOf(e)){for(;l!==null;){var o=Tr(l);if(o!==null&&ms(o),o=wo(e,t,n,r),o===null&&Mo(e,t,r,il,n),o===l)break;l=o}l!==null&&r.stopPropagation()}else Mo(e,t,r,null,n)}}var il=null;function wo(e,t,n,r){if(il=null,e=so(r),e=_n(e),e!==null)if(t=yn(e),t===null)e=null;else if(n=t.tag,n===13){if(e=ls(t),e!==null)return e;e=null}else if(n===3){if(t.stateNode.current.memoizedState.isDehydrated)return t.tag===3?t.stateNode.containerInfo:null;e=null}else t!==e&&(e=null);return il=e,null}function ks(e){switch(e){case"cancel":case"click":case"close":case"contextmenu":case"copy":case"cut":case"auxclick":case"dblclick":case"dragend":case"dragstart":case"drop":case"focusin":case"focusout":case"input":case"invalid":case"keydown":case"keypress":case"keyup":case"mousedown":case"mouseup":case"paste":case"pause":case"play":case"pointercancel":case"pointerdown":case"pointerup":case"ratechange":case"reset":case"resize":case"seeked":case"submit":case"touchcancel":case"touchend":case"touchstart":case"volumechange":case"change":case"selectionchange":case"textInput":case"compositionstart":case"compositionend":case"compositionupdate":case"beforeblur":case"afterblur":case"beforeinput":case"blur":case"fullscreenchange":case"focus":case"hashchange":case"popstate":case"select":case"selectstart":return 1;case"drag":case"dragenter":case"dragexit":case"dragleave":case"dragover":case"mousemove":case"mouseout":case"mouseover":case"pointermove":case"pointerout":case"pointerover":case"scroll":case"toggle":case"touchmove":case"wheel":case"mouseenter":case"mouseleave":case"pointerenter":case"pointerleave":return 4;case"message":switch(lc()){case po:return 1;case cs:return 4;case Jr:case oc:return 16;case fs:return 536870912;default:return 16}default:return 16}}var Xt=null,ko=null,sl=null;function Ss(){if(sl)return sl;var e,t=ko,n=t.length,r,l="value"in Xt?Xt.value:Xt.textContent,o=l.length;for(e=0;e<n&&t[e]===l[e];e++);var i=n-e;for(r=1;r<=i&&t[n-r]===l[o-r];r++);return sl=l.slice(e,1<r?1-r:void 0)}function al(e){var t=e.keyCode;return"charCode"in e?(e=e.charCode,e===0&&t===13&&(e=13)):e=t,e===10&&(e=13),32<=e||e===13?e:0}function ul(){return!0}function $s(){return!1}function lt(e){function t(n,r,l,o,i){this._reactName=n,this._targetInst=l,this.type=r,this.nativeEvent=o,this.target=i,this.currentTarget=null;for(var s in e)e.hasOwnProperty(s)&&(n=e[s],this[s]=n?n(o):o[s]);return this.isDefaultPrevented=(o.defaultPrevented!=null?o.defaultPrevented:o.returnValue===!1)?ul:$s,this.isPropagationStopped=$s,this}return w(t.prototype,{preventDefault:function(){this.defaultPrevented=!0;var n=this.nativeEvent;n&&(n.preventDefault?n.preventDefault():typeof n.returnValue!="unknown"&&(n.returnValue=!1),this.isDefaultPrevented=ul)},stopPropagation:function(){var n=this.nativeEvent;n&&(n.stopPropagation?n.stopPropagation():typeof n.cancelBubble!="unknown"&&(n.cancelBubble=!0),this.isPropagationStopped=ul)},persist:function(){},isPersistent:ul}),t}var Fn={eventPhase:0,bubbles:0,cancelable:0,timeStamp:function(e){return e.timeStamp||Date.now()},defaultPrevented:0,isTrusted:0},So=lt(Fn),yr=w({},Fn,{view:0,detail:0}),yc=lt(yr),$o,No,_r,cl=w({},yr,{screenX:0,screenY:0,clientX:0,clientY:0,pageX:0,pageY:0,ctrlKey:0,shiftKey:0,altKey:0,metaKey:0,getModifierState:Co,button:0,buttons:0,relatedTarget:function(e){return e.relatedTarget===void 0?e.fromElement===e.srcElement?e.toElement:e.fromElement:e.relatedTarget},movementX:function(e){return"movementX"in e?e.movementX:(e!==_r&&(_r&&e.type==="mousemove"?($o=e.screenX-_r.screenX,No=e.screenY-_r.screenY):No=$o=0,_r=e),$o)},movementY:function(e){return"movementY"in e?e.movementY:No}}),Ns=lt(cl),_c=w({},cl,{dataTransfer:0}),wc=lt(_c),kc=w({},yr,{relatedTarget:0}),Eo=lt(kc),Sc=w({},Fn,{animationName:0,elapsedTime:0,pseudoElement:0}),$c=lt(Sc),Nc=w({},Fn,{clipboardData:function(e){return"clipboardData"in e?e.clipboardData:window.clipboardData}}),Ec=lt(Nc),Cc=w({},Fn,{data:0}),Es=lt(Cc),xc={Esc:"Escape",Spacebar:" ",Left:"ArrowLeft",Up:"ArrowUp",Right:"ArrowRight",Down:"ArrowDown",Del:"Delete",Win:"OS",Menu:"ContextMenu",Apps:"ContextMenu",Scroll:"ScrollLock",MozPrintableKey:"Unidentified"},Pc={8:"Backspace",9:"Tab",12:"Clear",13:"Enter",16:"Shift",17:"Control",18:"Alt",19:"Pause",20:"CapsLock",27:"Escape",32:" ",33:"PageUp",34:"PageDown",35:"End",36:"Home",37:"ArrowLeft",38:"ArrowUp",39:"ArrowRight",40:"ArrowDown",45:"Insert",46:"Delete",112:"F1",113:"F2",114:"F3",115:"F4",116:"F5",117:"F6",118:"F7",119:"F8",120:"F9",121:"F10",122:"F11",123:"F12",144:"NumLock",145:"ScrollLock",224:"Meta"},Tc={Alt:"altKey",Control:"ctrlKey",Meta:"metaKey",Shift:"shiftKey"};function zc(e){var t=this.nativeEvent;return t.getModifierState?t.getModifierState(e):(e=Tc[e])?!!t[e]:!1}function Co(){return zc}var Rc=w({},yr,{key:function(e){if(e.key){var t=xc[e.key]||e.key;if(t!=="Unidentified")return t}return e.type==="keypress"?(e=al(e),e===13?"Enter":String.fromCharCode(e)):e.type==="keydown"||e.type==="keyup"?Pc[e.keyCode]||"Unidentified":""},code:0,location:0,ctrlKey:0,shiftKey:0,altKey:0,metaKey:0,repeat:0,locale:0,getModifierState:Co,charCode:function(e){return e.type==="keypress"?al(e):0},keyCode:function(e){return e.type==="keydown"||e.type==="keyup"?e.keyCode:0},which:function(e){return e.type==="keypress"?al(e):e.type==="keydown"||e.type==="keyup"?e.keyCode:0}}),Lc=lt(Rc),Ic=w({},cl,{pointerId:0,width:0,height:0,pressure:0,tangentialPressure:0,tiltX:0,tiltY:0,twist:0,pointerType:0,isPrimary:0}),Cs=lt(Ic),Dc=w({},yr,{touches:0,targetTouches:0,changedTouches:0,altKey:0,metaKey:0,ctrlKey:0,shiftKey:0,getModifierState:Co}),Oc=lt(Dc),Fc=w({},Fn,{propertyName:0,elapsedTime:0,pseudoElement:0}),Mc=lt(Fc),jc=w({},cl,{deltaX:function(e){return"deltaX"in e?e.deltaX:"wheelDeltaX"in e?-e.wheelDeltaX:0},deltaY:function(e){return"deltaY"in e?e.deltaY:"wheelDeltaY"in e?-e.wheelDeltaY:"wheelDelta"in e?-e.wheelDelta:0},deltaZ:0,deltaMode:0}),Ac=lt(jc),Uc=[9,13,27,32],xo=G&&"CompositionEvent"in window,wr=null;G&&"documentMode"in document&&(wr=document.documentMode);var Bc=G&&"TextEvent"in window&&!wr,xs=G&&(!xo||wr&&8<wr&&11>=wr),Ps=" ",Ts=!1;function zs(e,t){switch(e){case"keyup":return Uc.indexOf(t.keyCode)!==-1;case"keydown":return t.keyCode!==229;case"keypress":case"mousedown":case"focusout":return!0;default:return!1}}function Rs(e){return e=e.detail,typeof e=="object"&&"data"in e?e.data:null}var Mn=!1;function Wc(e,t){switch(e){case"compositionend":return Rs(t);case"keypress":return t.which!==32?null:(Ts=!0,Ps);case"textInput":return e=t.data,e===Ps&&Ts?null:e;default:return null}}function Qc(e,t){if(Mn)return e==="compositionend"||!xo&&zs(e,t)?(e=Ss(),sl=ko=Xt=null,Mn=!1,e):null;switch(e){case"paste":return null;case"keypress":if(!(t.ctrlKey||t.altKey||t.metaKey)||t.ctrlKey&&t.altKey){if(t.char&&1<t.char.length)return t.char;if(t.which)return String.fromCharCode(t.which)}return null;case"compositionend":return xs&&t.locale!=="ko"?null:t.data;default:return null}}var Vc={color:!0,date:!0,datetime:!0,"datetime-local":!0,email:!0,month:!0,number:!0,password:!0,range:!0,search:!0,tel:!0,text:!0,time:!0,url:!0,week:!0};function Ls(e){var t=e&&e.nodeName&&e.nodeName.toLowerCase();return t==="input"?!!Vc[e.type]:t==="textarea"}function Is(e,t,n,r){bi(r),t=vl(t,"onChange"),0<t.length&&(n=new So("onChange","change",null,n,r),e.push({event:n,listeners:t}))}var kr=null,Sr=null;function Hc(e){Zs(e,0)}function fl(e){var t=Wn(e);if(dn(t))return e}function Kc(e,t){if(e==="change")return t}var Ds=!1;if(G){var Po;if(G){var To="oninput"in document;if(!To){var Os=document.createElement("div");Os.setAttribute("oninput","return;"),To=typeof Os.oninput=="function"}Po=To}else Po=!1;Ds=Po&&(!document.documentMode||9<document.documentMode)}function Fs(){kr&&(kr.detachEvent("onpropertychange",Ms),Sr=kr=null)}function Ms(e){if(e.propertyName==="value"&&fl(Sr)){var t=[];Is(t,Sr,e,so(e)),rs(Hc,t)}}function Yc(e,t,n){e==="focusin"?(Fs(),kr=t,Sr=n,kr.attachEvent("onpropertychange",Ms)):e==="focusout"&&Fs()}function Gc(e){if(e==="selectionchange"||e==="keyup"||e==="keydown")return fl(Sr)}function Xc(e,t){if(e==="click")return fl(t)}function qc(e,t){if(e==="input"||e==="change")return fl(t)}function Zc(e,t){return e===t&&(e!==0||1/e===1/t)||e!==e&&t!==t}var wt=typeof Object.is=="function"?Object.is:Zc;function $r(e,t){if(wt(e,t))return!0;if(typeof e!="object"||e===null||typeof t!="object"||t===null)return!1;var n=Object.keys(e),r=Object.keys(t);if(n.length!==r.length)return!1;for(r=0;r<n.length;r++){var l=n[r];if(!E.call(t,l)||!wt(e[l],t[l]))return!1}return!0}function js(e){for(;e&&e.firstChild;)e=e.firstChild;return e}function As(e,t){var n=js(e);e=0;for(var r;n;){if(n.nodeType===3){if(r=e+n.textContent.length,e<=t&&r>=t)return{node:n,offset:t-e};e=r}e:{for(;n;){if(n.nextSibling){n=n.nextSibling;break e}n=n.parentNode}n=void 0}n=js(n)}}function Us(e,t){return e&&t?e===t?!0:e&&e.nodeType===3?!1:t&&t.nodeType===3?Us(e,t.parentNode):"contains"in e?e.contains(t):e.compareDocumentPosition?!!(e.compareDocumentPosition(t)&16):!1:!1}function Bs(){for(var e=window,t=pn();t instanceof e.HTMLIFrameElement;){try{var n=typeof t.contentWindow.location.href=="string"}catch{n=!1}if(n)e=t.contentWindow;else break;t=pn(e.document)}return t}function zo(e){var t=e&&e.nodeName&&e.nodeName.toLowerCase();return t&&(t==="input"&&(e.type==="text"||e.type==="search"||e.type==="tel"||e.type==="url"||e.type==="password")||t==="textarea"||e.contentEditable==="true")}function Jc(e){var t=Bs(),n=e.focusedElem,r=e.selectionRange;if(t!==n&&n&&n.ownerDocument&&Us(n.ownerDocument.documentElement,n)){if(r!==null&&zo(n)){if(t=r.start,e=r.end,e===void 0&&(e=t),"selectionStart"in n)n.selectionStart=t,n.selectionEnd=Math.min(e,n.value.length);else if(e=(t=n.ownerDocument||document)&&t.defaultView||window,e.getSelection){e=e.getSelection();var l=n.textContent.length,o=Math.min(r.start,l);r=r.end===void 0?o:Math.min(r.end,l),!e.extend&&o>r&&(l=r,r=o,o=l),l=As(n,o);var i=As(n,r);l&&i&&(e.rangeCount!==1||e.anchorNode!==l.node||e.anchorOffset!==l.offset||e.focusNode!==i.node||e.focusOffset!==i.offset)&&(t=t.createRange(),t.setStart(l.node,l.offset),e.removeAllRanges(),o>r?(e.addRange(t),e.extend(i.node,i.offset)):(t.setEnd(i.node,i.offset),e.addRange(t)))}}for(t=[],e=n;e=e.parentNode;)e.nodeType===1&&t.push({element:e,left:e.scrollLeft,top:e.scrollTop});for(typeof n.focus=="function"&&n.focus(),n=0;n<t.length;n++)e=t[n],e.element.scrollLeft=e.left,e.element.scrollTop=e.top}}var bc=G&&"documentMode"in document&&11>=document.documentMode,jn=null,Ro=null,Nr=null,Lo=!1;function Ws(e,t,n){var r=n.window===n?n.document:n.nodeType===9?n:n.ownerDocument;Lo||jn==null||jn!==pn(r)||(r=jn,"selectionStart"in r&&zo(r)?r={start:r.selectionStart,end:r.selectionEnd}:(r=(r.ownerDocument&&r.ownerDocument.defaultView||window).getSelection(),r={anchorNode:r.anchorNode,anchorOffset:r.anchorOffset,focusNode:r.focusNode,focusOffset:r.focusOffset}),Nr&&$r(Nr,r)||(Nr=r,r=vl(Ro,"onSelect"),0<r.length&&(t=new So("onSelect","select",null,t,n),e.push({event:t,listeners:r}),t.target=jn)))}function dl(e,t){var n={};return n[e.toLowerCase()]=t.toLowerCase(),n["Webkit"+e]="webkit"+t,n["Moz"+e]="moz"+t,n}var An={animationend:dl("Animation","AnimationEnd"),animationiteration:dl("Animation","AnimationIteration"),animationstart:dl("Animation","AnimationStart"),transitionend:dl("Transition","TransitionEnd")},Io={},Qs={};G&&(Qs=document.createElement("div").style,"AnimationEvent"in window||(delete An.animationend.animation,delete An.animationiteration.animation,delete An.animationstart.animation),"TransitionEvent"in window||delete An.transitionend.transition);function pl(e){if(Io[e])return Io[e];if(!An[e])return e;var t=An[e],n;for(n in t)if(t.hasOwnProperty(n)&&n in Qs)return Io[e]=t[n];return e}var Vs=pl("animationend"),Hs=pl("animationiteration"),Ks=pl("animationstart"),Ys=pl("transitionend"),Gs=new Map,Xs="abort auxClick cancel canPlay canPlayThrough click close contextMenu copy cut drag dragEnd dragEnter dragExit dragLeave dragOver dragStart drop durationChange emptied encrypted ended error gotPointerCapture input invalid keyDown keyPress keyUp load loadedData loadedMetadata loadStart lostPointerCapture mouseDown mouseMove mouseOut mouseOver mouseUp paste pause play playing pointerCancel pointerDown pointerMove pointerOut pointerOver pointerUp progress rateChange reset resize seeked seeking stalled submit suspend timeUpdate touchCancel touchEnd touchStart volumeChange scroll toggle touchMove waiting wheel".split(" ");function qt(e,t){Gs.set(e,t),x(t,[e])}for(var Do=0;Do<Xs.length;Do++){var Oo=Xs[Do],ef=Oo.toLowerCase(),tf=Oo[0].toUpperCase()+Oo.slice(1);qt(ef,"on"+tf)}qt(Vs,"onAnimationEnd"),qt(Hs,"onAnimationIteration"),qt(Ks,"onAnimationStart"),qt("dblclick","onDoubleClick"),qt("focusin","onFocus"),qt("focusout","onBlur"),qt(Ys,"onTransitionEnd"),M("onMouseEnter",["mouseout","mouseover"]),M("onMouseLeave",["mouseout","mouseover"]),M("onPointerEnter",["pointerout","pointerover"]),M("onPointerLeave",["pointerout","pointerover"]),x("onChange","change click focusin focusout input keydown keyup selectionchange".split(" ")),x("onSelect","focusout contextmenu dragend focusin keydown keyup mousedown mouseup selectionchange".split(" ")),x("onBeforeInput",["compositionend","keypress","textInput","paste"]),x("onCompositionEnd","compositionend focusout keydown keypress keyup mousedown".split(" ")),x("onCompositionStart","compositionstart focusout keydown keypress keyup mousedown".split(" ")),x("onCompositionUpdate","compositionupdate focusout keydown keypress keyup mousedown".split(" "));var Er="abort canplay canplaythrough durationchange emptied encrypted ended error loadeddata loadedmetadata loadstart pause play playing progress ratechange resize seeked seeking stalled suspend timeupdate volumechange waiting".split(" "),nf=new Set("cancel close invalid load scroll toggle".split(" ").concat(Er));function qs(e,t,n){var r=e.type||"unknown-event";e.currentTarget=n,ec(r,t,void 0,e),e.currentTarget=null}function Zs(e,t){t=(t&4)!==0;for(var n=0;n<e.length;n++){var r=e[n],l=r.event;r=r.listeners;e:{var o=void 0;if(t)for(var i=r.length-1;0<=i;i--){var s=r[i],c=s.instance,g=s.currentTarget;if(s=s.listener,c!==o&&l.isPropagationStopped())break e;qs(l,s,g),o=c}else for(i=0;i<r.length;i++){if(s=r[i],c=s.instance,g=s.currentTarget,s=s.listener,c!==o&&l.isPropagationStopped())break e;qs(l,s,g),o=c}}}if(Zr)throw e=fo,Zr=!1,fo=null,e}function fe(e,t){var n=t[Qo];n===void 0&&(n=t[Qo]=new Set);var r=e+"__bubble";n.has(r)||(Js(t,e,2,!1),n.add(r))}function Fo(e,t,n){var r=0;t&&(r|=4),Js(n,e,r,t)}var ml="_reactListening"+Math.random().toString(36).slice(2);function Cr(e){if(!e[ml]){e[ml]=!0,T.forEach(function(n){n!=="selectionchange"&&(nf.has(n)||Fo(n,!1,e),Fo(n,!0,e))});var t=e.nodeType===9?e:e.ownerDocument;t===null||t[ml]||(t[ml]=!0,Fo("selectionchange",!1,t))}}function Js(e,t,n,r){switch(ks(t)){case 1:var l=hc;break;case 4:l=gc;break;default:l=_o}n=l.bind(null,t,n,e),l=void 0,!co||t!=="touchstart"&&t!=="touchmove"&&t!=="wheel"||(l=!0),r?l!==void 0?e.addEventListener(t,n,{capture:!0,passive:l}):e.addEventListener(t,n,!0):l!==void 0?e.addEventListener(t,n,{passive:l}):e.addEventListener(t,n,!1)}function Mo(e,t,n,r,l){var o=r;if((t&1)===0&&(t&2)===0&&r!==null)e:for(;;){if(r===null)return;var i=r.tag;if(i===3||i===4){var s=r.stateNode.containerInfo;if(s===l||s.nodeType===8&&s.parentNode===l)break;if(i===4)for(i=r.return;i!==null;){var c=i.tag;if((c===3||c===4)&&(c=i.stateNode.containerInfo,c===l||c.nodeType===8&&c.parentNode===l))return;i=i.return}for(;s!==null;){if(i=_n(s),i===null)return;if(c=i.tag,c===5||c===6){r=o=i;continue e}s=s.parentNode}}r=r.return}rs(function(){var g=o,S=so(n),$=[];e:{var k=Gs.get(e);if(k!==void 0){var z=So,O=e;switch(e){case"keypress":if(al(n)===0)break e;case"keydown":case"keyup":z=Lc;break;case"focusin":O="focus",z=Eo;break;case"focusout":O="blur",z=Eo;break;case"beforeblur":case"afterblur":z=Eo;break;case"click":if(n.button===2)break e;case"auxclick":case"dblclick":case"mousedown":case"mousemove":case"mouseup":case"mouseout":case"mouseover":case"contextmenu":z=Ns;break;case"drag":case"dragend":case"dragenter":case"dragexit":case"dragleave":case"dragover":case"dragstart":case"drop":z=wc;break;case"touchcancel":case"touchend":case"touchmove":case"touchstart":z=Oc;break;case Vs:case Hs:case Ks:z=$c;break;case Ys:z=Mc;break;case"scroll":z=yc;break;case"wheel":z=Ac;break;case"copy":case"cut":case"paste":z=Ec;break;case"gotpointercapture":case"lostpointercapture":case"pointercancel":case"pointerdown":case"pointermove":case"pointerout":case"pointerover":case"pointerup":z=Cs}var F=(t&4)!==0,Se=!F&&e==="scroll",v=F?k!==null?k+"Capture":null:k;F=[];for(var p=g,h;p!==null;){h=p;var N=h.stateNode;if(h.tag===5&&N!==null&&(h=N,v!==null&&(N=ar(p,v),N!=null&&F.push(xr(p,N,h)))),Se)break;p=p.return}0<F.length&&(k=new z(k,O,null,n,S),$.push({event:k,listeners:F}))}}if((t&7)===0){e:{if(k=e==="mouseover"||e==="pointerover",z=e==="mouseout"||e==="pointerout",k&&n!==io&&(O=n.relatedTarget||n.fromElement)&&(_n(O)||O[It]))break e;if((z||k)&&(k=S.window===S?S:(k=S.ownerDocument)?k.defaultView||k.parentWindow:window,z?(O=n.relatedTarget||n.toElement,z=g,O=O?_n(O):null,O!==null&&(Se=yn(O),O!==Se||O.tag!==5&&O.tag!==6)&&(O=null)):(z=null,O=g),z!==O)){if(F=Ns,N="onMouseLeave",v="onMouseEnter",p="mouse",(e==="pointerout"||e==="pointerover")&&(F=Cs,N="onPointerLeave",v="onPointerEnter",p="pointer"),Se=z==null?k:Wn(z),h=O==null?k:Wn(O),k=new F(N,p+"leave",z,n,S),k.target=Se,k.relatedTarget=h,N=null,_n(S)===g&&(F=new F(v,p+"enter",O,n,S),F.target=h,F.relatedTarget=Se,N=F),Se=N,z&&O)t:{for(F=z,v=O,p=0,h=F;h;h=Un(h))p++;for(h=0,N=v;N;N=Un(N))h++;for(;0<p-h;)F=Un(F),p--;for(;0<h-p;)v=Un(v),h--;for(;p--;){if(F===v||v!==null&&F===v.alternate)break t;F=Un(F),v=Un(v)}F=null}else F=null;z!==null&&bs($,k,z,F,!1),O!==null&&Se!==null&&bs($,Se,O,F,!0)}}e:{if(k=g?Wn(g):window,z=k.nodeName&&k.nodeName.toLowerCase(),z==="select"||z==="input"&&k.type==="file")var j=Kc;else if(Ls(k))if(Ds)j=qc;else{j=Gc;var B=Yc}else(z=k.nodeName)&&z.toLowerCase()==="input"&&(k.type==="checkbox"||k.type==="radio")&&(j=Xc);if(j&&(j=j(e,g))){Is($,j,n,S);break e}B&&B(e,k,g),e==="focusout"&&(B=k._wrapperState)&&B.controlled&&k.type==="number"&&ir(k,"number",k.value)}switch(B=g?Wn(g):window,e){case"focusin":(Ls(B)||B.contentEditable==="true")&&(jn=B,Ro=g,Nr=null);break;case"focusout":Nr=Ro=jn=null;break;case"mousedown":Lo=!0;break;case"contextmenu":case"mouseup":case"dragend":Lo=!1,Ws($,n,S);break;case"selectionchange":if(bc)break;case"keydown":case"keyup":Ws($,n,S)}var W;if(xo)e:{switch(e){case"compositionstart":var V="onCompositionStart";break e;case"compositionend":V="onCompositionEnd";break e;case"compositionupdate":V="onCompositionUpdate";break e}V=void 0}else Mn?zs(e,n)&&(V="onCompositionEnd"):e==="keydown"&&n.keyCode===229&&(V="onCompositionStart");V&&(xs&&n.locale!=="ko"&&(Mn||V!=="onCompositionStart"?V==="onCompositionEnd"&&Mn&&(W=Ss()):(Xt=S,ko="value"in Xt?Xt.value:Xt.textContent,Mn=!0)),B=vl(g,V),0<B.length&&(V=new Es(V,e,null,n,S),$.push({event:V,listeners:B}),W?V.data=W:(W=Rs(n),W!==null&&(V.data=W)))),(W=Bc?Wc(e,n):Qc(e,n))&&(g=vl(g,"onBeforeInput"),0<g.length&&(S=new Es("onBeforeInput","beforeinput",null,n,S),$.push({event:S,listeners:g}),S.data=W))}Zs($,t)})}function xr(e,t,n){return{instance:e,listener:t,currentTarget:n}}function vl(e,t){for(var n=t+"Capture",r=[];e!==null;){var l=e,o=l.stateNode;l.tag===5&&o!==null&&(l=o,o=ar(e,n),o!=null&&r.unshift(xr(e,o,l)),o=ar(e,t),o!=null&&r.push(xr(e,o,l))),e=e.return}return r}function Un(e){if(e===null)return null;do e=e.return;while(e&&e.tag!==5);return e||null}function bs(e,t,n,r,l){for(var o=t._reactName,i=[];n!==null&&n!==r;){var s=n,c=s.alternate,g=s.stateNode;if(c!==null&&c===r)break;s.tag===5&&g!==null&&(s=g,l?(c=ar(n,o),c!=null&&i.unshift(xr(n,c,s))):l||(c=ar(n,o),c!=null&&i.push(xr(n,c,s)))),n=n.return}i.length!==0&&e.push({event:t,listeners:i})}var rf=/\r\n?/g,lf=/\u0000|\uFFFD/g;function ea(e){return(typeof e=="string"?e:""+e).replace(rf,`
`).replace(lf,"")}function hl(e,t,n){if(t=ea(t),ea(e)!==t&&n)throw Error(m(425))}function gl(){}var jo=null,Ao=null;function Uo(e,t){return e==="textarea"||e==="noscript"||typeof t.children=="string"||typeof t.children=="number"||typeof t.dangerouslySetInnerHTML=="object"&&t.dangerouslySetInnerHTML!==null&&t.dangerouslySetInnerHTML.__html!=null}var Bo=typeof setTimeout=="function"?setTimeout:void 0,of=typeof clearTimeout=="function"?clearTimeout:void 0,ta=typeof Promise=="function"?Promise:void 0,sf=typeof queueMicrotask=="function"?queueMicrotask:typeof ta<"u"?function(e){return ta.resolve(null).then(e).catch(af)}:Bo;function af(e){setTimeout(function(){throw e})}function Wo(e,t){var n=t,r=0;do{var l=n.nextSibling;if(e.removeChild(n),l&&l.nodeType===8)if(n=l.data,n==="/$"){if(r===0){e.removeChild(l),gr(t);return}r--}else n!=="$"&&n!=="$?"&&n!=="$!"||r++;n=l}while(n);gr(t)}function Zt(e){for(;e!=null;e=e.nextSibling){var t=e.nodeType;if(t===1||t===3)break;if(t===8){if(t=e.data,t==="$"||t==="$!"||t==="$?")break;if(t==="/$")return null}}return e}function na(e){e=e.previousSibling;for(var t=0;e;){if(e.nodeType===8){var n=e.data;if(n==="$"||n==="$!"||n==="$?"){if(t===0)return e;t--}else n==="/$"&&t++}e=e.previousSibling}return null}var Bn=Math.random().toString(36).slice(2),Pt="__reactFiber$"+Bn,Pr="__reactProps$"+Bn,It="__reactContainer$"+Bn,Qo="__reactEvents$"+Bn,uf="__reactListeners$"+Bn,cf="__reactHandles$"+Bn;function _n(e){var t=e[Pt];if(t)return t;for(var n=e.parentNode;n;){if(t=n[It]||n[Pt]){if(n=t.alternate,t.child!==null||n!==null&&n.child!==null)for(e=na(e);e!==null;){if(n=e[Pt])return n;e=na(e)}return t}e=n,n=e.parentNode}return null}function Tr(e){return e=e[Pt]||e[It],!e||e.tag!==5&&e.tag!==6&&e.tag!==13&&e.tag!==3?null:e}function Wn(e){if(e.tag===5||e.tag===6)return e.stateNode;throw Error(m(33))}function yl(e){return e[Pr]||null}var Vo=[],Qn=-1;function Jt(e){return{current:e}}function de(e){0>Qn||(e.current=Vo[Qn],Vo[Qn]=null,Qn--)}function ce(e,t){Qn++,Vo[Qn]=e.current,e.current=t}var bt={},We=Jt(bt),Xe=Jt(!1),wn=bt;function Vn(e,t){var n=e.type.contextTypes;if(!n)return bt;var r=e.stateNode;if(r&&r.__reactInternalMemoizedUnmaskedChildContext===t)return r.__reactInternalMemoizedMaskedChildContext;var l={},o;for(o in n)l[o]=t[o];return r&&(e=e.stateNode,e.__reactInternalMemoizedUnmaskedChildContext=t,e.__reactInternalMemoizedMaskedChildContext=l),l}function qe(e){return e=e.childContextTypes,e!=null}function _l(){de(Xe),de(We)}function ra(e,t,n){if(We.current!==bt)throw Error(m(168));ce(We,t),ce(Xe,n)}function la(e,t,n){var r=e.stateNode;if(t=t.childContextTypes,typeof r.getChildContext!="function")return n;r=r.getChildContext();for(var l in r)if(!(l in t))throw Error(m(108,J(e)||"Unknown",l));return w({},n,r)}function wl(e){return e=(e=e.stateNode)&&e.__reactInternalMemoizedMergedChildContext||bt,wn=We.current,ce(We,e),ce(Xe,Xe.current),!0}function oa(e,t,n){var r=e.stateNode;if(!r)throw Error(m(169));n?(e=la(e,t,wn),r.__reactInternalMemoizedMergedChildContext=e,de(Xe),de(We),ce(We,e)):de(Xe),ce(Xe,n)}var Dt=null,kl=!1,Ho=!1;function ia(e){Dt===null?Dt=[e]:Dt.push(e)}function ff(e){kl=!0,ia(e)}function en(){if(!Ho&&Dt!==null){Ho=!0;var e=0,t=ie;try{var n=Dt;for(ie=1;e<n.length;e++){var r=n[e];do r=r(!0);while(r!==null)}Dt=null,kl=!1}catch(l){throw Dt!==null&&(Dt=Dt.slice(e+1)),as(po,en),l}finally{ie=t,Ho=!1}}return null}var Hn=[],Kn=0,Sl=null,$l=0,ct=[],ft=0,kn=null,Ot=1,Ft="";function Sn(e,t){Hn[Kn++]=$l,Hn[Kn++]=Sl,Sl=e,$l=t}function sa(e,t,n){ct[ft++]=Ot,ct[ft++]=Ft,ct[ft++]=kn,kn=e;var r=Ot;e=Ft;var l=32-_t(r)-1;r&=~(1<<l),n+=1;var o=32-_t(t)+l;if(30<o){var i=l-l%5;o=(r&(1<<i)-1).toString(32),r>>=i,l-=i,Ot=1<<32-_t(t)+l|n<<l|r,Ft=o+e}else Ot=1<<o|n<<l|r,Ft=e}function Ko(e){e.return!==null&&(Sn(e,1),sa(e,1,0))}function Yo(e){for(;e===Sl;)Sl=Hn[--Kn],Hn[Kn]=null,$l=Hn[--Kn],Hn[Kn]=null;for(;e===kn;)kn=ct[--ft],ct[ft]=null,Ft=ct[--ft],ct[ft]=null,Ot=ct[--ft],ct[ft]=null}var ot=null,it=null,ve=!1,kt=null;function aa(e,t){var n=vt(5,null,null,0);n.elementType="DELETED",n.stateNode=t,n.return=e,t=e.deletions,t===null?(e.deletions=[n],e.flags|=16):t.push(n)}function ua(e,t){switch(e.tag){case 5:var n=e.type;return t=t.nodeType!==1||n.toLowerCase()!==t.nodeName.toLowerCase()?null:t,t!==null?(e.stateNode=t,ot=e,it=Zt(t.firstChild),!0):!1;case 6:return t=e.pendingProps===""||t.nodeType!==3?null:t,t!==null?(e.stateNode=t,ot=e,it=null,!0):!1;case 13:return t=t.nodeType!==8?null:t,t!==null?(n=kn!==null?{id:Ot,overflow:Ft}:null,e.memoizedState={dehydrated:t,treeContext:n,retryLane:1073741824},n=vt(18,null,null,0),n.stateNode=t,n.return=e,e.child=n,ot=e,it=null,!0):!1;default:return!1}}function Go(e){return(e.mode&1)!==0&&(e.flags&128)===0}function Xo(e){if(ve){var t=it;if(t){var n=t;if(!ua(e,t)){if(Go(e))throw Error(m(418));t=Zt(n.nextSibling);var r=ot;t&&ua(e,t)?aa(r,n):(e.flags=e.flags&-4097|2,ve=!1,ot=e)}}else{if(Go(e))throw Error(m(418));e.flags=e.flags&-4097|2,ve=!1,ot=e}}}function ca(e){for(e=e.return;e!==null&&e.tag!==5&&e.tag!==3&&e.tag!==13;)e=e.return;ot=e}function Nl(e){if(e!==ot)return!1;if(!ve)return ca(e),ve=!0,!1;var t;if((t=e.tag!==3)&&!(t=e.tag!==5)&&(t=e.type,t=t!=="head"&&t!=="body"&&!Uo(e.type,e.memoizedProps)),t&&(t=it)){if(Go(e))throw fa(),Error(m(418));for(;t;)aa(e,t),t=Zt(t.nextSibling)}if(ca(e),e.tag===13){if(e=e.memoizedState,e=e!==null?e.dehydrated:null,!e)throw Error(m(317));e:{for(e=e.nextSibling,t=0;e;){if(e.nodeType===8){var n=e.data;if(n==="/$"){if(t===0){it=Zt(e.nextSibling);break e}t--}else n!=="$"&&n!=="$!"&&n!=="$?"||t++}e=e.nextSibling}it=null}}else it=ot?Zt(e.stateNode.nextSibling):null;return!0}function fa(){for(var e=it;e;)e=Zt(e.nextSibling)}function Yn(){it=ot=null,ve=!1}function qo(e){kt===null?kt=[e]:kt.push(e)}var df=we.ReactCurrentBatchConfig;function zr(e,t,n){if(e=n.ref,e!==null&&typeof e!="function"&&typeof e!="object"){if(n._owner){if(n=n._owner,n){if(n.tag!==1)throw Error(m(309));var r=n.stateNode}if(!r)throw Error(m(147,e));var l=r,o=""+e;return t!==null&&t.ref!==null&&typeof t.ref=="function"&&t.ref._stringRef===o?t.ref:(t=function(i){var s=l.refs;i===null?delete s[o]:s[o]=i},t._stringRef=o,t)}if(typeof e!="string")throw Error(m(284));if(!n._owner)throw Error(m(290,e))}return e}function El(e,t){throw e=Object.prototype.toString.call(t),Error(m(31,e==="[object Object]"?"object with keys {"+Object.keys(t).join(", ")+"}":e))}function da(e){var t=e._init;return t(e._payload)}function pa(e){function t(v,p){if(e){var h=v.deletions;h===null?(v.deletions=[p],v.flags|=16):h.push(p)}}function n(v,p){if(!e)return null;for(;p!==null;)t(v,p),p=p.sibling;return null}function r(v,p){for(v=new Map;p!==null;)p.key!==null?v.set(p.key,p):v.set(p.index,p),p=p.sibling;return v}function l(v,p){return v=un(v,p),v.index=0,v.sibling=null,v}function o(v,p,h){return v.index=h,e?(h=v.alternate,h!==null?(h=h.index,h<p?(v.flags|=2,p):h):(v.flags|=2,p)):(v.flags|=1048576,p)}function i(v){return e&&v.alternate===null&&(v.flags|=2),v}function s(v,p,h,N){return p===null||p.tag!==6?(p=Bi(h,v.mode,N),p.return=v,p):(p=l(p,h),p.return=v,p)}function c(v,p,h,N){var j=h.type;return j===Fe?S(v,p,h.props.children,N,h.key):p!==null&&(p.elementType===j||typeof j=="object"&&j!==null&&j.$$typeof===U&&da(j)===p.type)?(N=l(p,h.props),N.ref=zr(v,p,h),N.return=v,N):(N=Xl(h.type,h.key,h.props,null,v.mode,N),N.ref=zr(v,p,h),N.return=v,N)}function g(v,p,h,N){return p===null||p.tag!==4||p.stateNode.containerInfo!==h.containerInfo||p.stateNode.implementation!==h.implementation?(p=Wi(h,v.mode,N),p.return=v,p):(p=l(p,h.children||[]),p.return=v,p)}function S(v,p,h,N,j){return p===null||p.tag!==7?(p=zn(h,v.mode,N,j),p.return=v,p):(p=l(p,h),p.return=v,p)}function $(v,p,h){if(typeof p=="string"&&p!==""||typeof p=="number")return p=Bi(""+p,v.mode,h),p.return=v,p;if(typeof p=="object"&&p!==null){switch(p.$$typeof){case rt:return h=Xl(p.type,p.key,p.props,null,v.mode,h),h.ref=zr(v,null,p),h.return=v,h;case Oe:return p=Wi(p,v.mode,h),p.return=v,p;case U:var N=p._init;return $(v,N(p._payload),h)}if(vn(p)||A(p))return p=zn(p,v.mode,h,null),p.return=v,p;El(v,p)}return null}function k(v,p,h,N){var j=p!==null?p.key:null;if(typeof h=="string"&&h!==""||typeof h=="number")return j!==null?null:s(v,p,""+h,N);if(typeof h=="object"&&h!==null){switch(h.$$typeof){case rt:return h.key===j?c(v,p,h,N):null;case Oe:return h.key===j?g(v,p,h,N):null;case U:return j=h._init,k(v,p,j(h._payload),N)}if(vn(h)||A(h))return j!==null?null:S(v,p,h,N,null);El(v,h)}return null}function z(v,p,h,N,j){if(typeof N=="string"&&N!==""||typeof N=="number")return v=v.get(h)||null,s(p,v,""+N,j);if(typeof N=="object"&&N!==null){switch(N.$$typeof){case rt:return v=v.get(N.key===null?h:N.key)||null,c(p,v,N,j);case Oe:return v=v.get(N.key===null?h:N.key)||null,g(p,v,N,j);case U:var B=N._init;return z(v,p,h,B(N._payload),j)}if(vn(N)||A(N))return v=v.get(h)||null,S(p,v,N,j,null);El(p,N)}return null}function O(v,p,h,N){for(var j=null,B=null,W=p,V=p=0,Ie=null;W!==null&&V<h.length;V++){W.index>V?(Ie=W,W=null):Ie=W.sibling;var te=k(v,W,h[V],N);if(te===null){W===null&&(W=Ie);break}e&&W&&te.alternate===null&&t(v,W),p=o(te,p,V),B===null?j=te:B.sibling=te,B=te,W=Ie}if(V===h.length)return n(v,W),ve&&Sn(v,V),j;if(W===null){for(;V<h.length;V++)W=$(v,h[V],N),W!==null&&(p=o(W,p,V),B===null?j=W:B.sibling=W,B=W);return ve&&Sn(v,V),j}for(W=r(v,W);V<h.length;V++)Ie=z(W,v,V,h[V],N),Ie!==null&&(e&&Ie.alternate!==null&&W.delete(Ie.key===null?V:Ie.key),p=o(Ie,p,V),B===null?j=Ie:B.sibling=Ie,B=Ie);return e&&W.forEach(function(cn){return t(v,cn)}),ve&&Sn(v,V),j}function F(v,p,h,N){var j=A(h);if(typeof j!="function")throw Error(m(150));if(h=j.call(h),h==null)throw Error(m(151));for(var B=j=null,W=p,V=p=0,Ie=null,te=h.next();W!==null&&!te.done;V++,te=h.next()){W.index>V?(Ie=W,W=null):Ie=W.sibling;var cn=k(v,W,te.value,N);if(cn===null){W===null&&(W=Ie);break}e&&W&&cn.alternate===null&&t(v,W),p=o(cn,p,V),B===null?j=cn:B.sibling=cn,B=cn,W=Ie}if(te.done)return n(v,W),ve&&Sn(v,V),j;if(W===null){for(;!te.done;V++,te=h.next())te=$(v,te.value,N),te!==null&&(p=o(te,p,V),B===null?j=te:B.sibling=te,B=te);return ve&&Sn(v,V),j}for(W=r(v,W);!te.done;V++,te=h.next())te=z(W,v,V,te.value,N),te!==null&&(e&&te.alternate!==null&&W.delete(te.key===null?V:te.key),p=o(te,p,V),B===null?j=te:B.sibling=te,B=te);return e&&W.forEach(function(Hf){return t(v,Hf)}),ve&&Sn(v,V),j}function Se(v,p,h,N){if(typeof h=="object"&&h!==null&&h.type===Fe&&h.key===null&&(h=h.props.children),typeof h=="object"&&h!==null){switch(h.$$typeof){case rt:e:{for(var j=h.key,B=p;B!==null;){if(B.key===j){if(j=h.type,j===Fe){if(B.tag===7){n(v,B.sibling),p=l(B,h.props.children),p.return=v,v=p;break e}}else if(B.elementType===j||typeof j=="object"&&j!==null&&j.$$typeof===U&&da(j)===B.type){n(v,B.sibling),p=l(B,h.props),p.ref=zr(v,B,h),p.return=v,v=p;break e}n(v,B);break}else t(v,B);B=B.sibling}h.type===Fe?(p=zn(h.props.children,v.mode,N,h.key),p.return=v,v=p):(N=Xl(h.type,h.key,h.props,null,v.mode,N),N.ref=zr(v,p,h),N.return=v,v=N)}return i(v);case Oe:e:{for(B=h.key;p!==null;){if(p.key===B)if(p.tag===4&&p.stateNode.containerInfo===h.containerInfo&&p.stateNode.implementation===h.implementation){n(v,p.sibling),p=l(p,h.children||[]),p.return=v,v=p;break e}else{n(v,p);break}else t(v,p);p=p.sibling}p=Wi(h,v.mode,N),p.return=v,v=p}return i(v);case U:return B=h._init,Se(v,p,B(h._payload),N)}if(vn(h))return O(v,p,h,N);if(A(h))return F(v,p,h,N);El(v,h)}return typeof h=="string"&&h!==""||typeof h=="number"?(h=""+h,p!==null&&p.tag===6?(n(v,p.sibling),p=l(p,h),p.return=v,v=p):(n(v,p),p=Bi(h,v.mode,N),p.return=v,v=p),i(v)):n(v,p)}return Se}var Gn=pa(!0),ma=pa(!1),Cl=Jt(null),xl=null,Xn=null,Zo=null;function Jo(){Zo=Xn=xl=null}function bo(e){var t=Cl.current;de(Cl),e._currentValue=t}function ei(e,t,n){for(;e!==null;){var r=e.alternate;if((e.childLanes&t)!==t?(e.childLanes|=t,r!==null&&(r.childLanes|=t)):r!==null&&(r.childLanes&t)!==t&&(r.childLanes|=t),e===n)break;e=e.return}}function qn(e,t){xl=e,Zo=Xn=null,e=e.dependencies,e!==null&&e.firstContext!==null&&((e.lanes&t)!==0&&(Ze=!0),e.firstContext=null)}function dt(e){var t=e._currentValue;if(Zo!==e)if(e={context:e,memoizedValue:t,next:null},Xn===null){if(xl===null)throw Error(m(308));Xn=e,xl.dependencies={lanes:0,firstContext:e}}else Xn=Xn.next=e;return t}var $n=null;function ti(e){$n===null?$n=[e]:$n.push(e)}function va(e,t,n,r){var l=t.interleaved;return l===null?(n.next=n,ti(t)):(n.next=l.next,l.next=n),t.interleaved=n,Mt(e,r)}function Mt(e,t){e.lanes|=t;var n=e.alternate;for(n!==null&&(n.lanes|=t),n=e,e=e.return;e!==null;)e.childLanes|=t,n=e.alternate,n!==null&&(n.childLanes|=t),n=e,e=e.return;return n.tag===3?n.stateNode:null}var tn=!1;function ni(e){e.updateQueue={baseState:e.memoizedState,firstBaseUpdate:null,lastBaseUpdate:null,shared:{pending:null,interleaved:null,lanes:0},effects:null}}function ha(e,t){e=e.updateQueue,t.updateQueue===e&&(t.updateQueue={baseState:e.baseState,firstBaseUpdate:e.firstBaseUpdate,lastBaseUpdate:e.lastBaseUpdate,shared:e.shared,effects:e.effects})}function jt(e,t){return{eventTime:e,lane:t,tag:0,payload:null,callback:null,next:null}}function nn(e,t,n){var r=e.updateQueue;if(r===null)return null;if(r=r.shared,(ee&2)!==0){var l=r.pending;return l===null?t.next=t:(t.next=l.next,l.next=t),r.pending=t,Mt(e,n)}return l=r.interleaved,l===null?(t.next=t,ti(r)):(t.next=l.next,l.next=t),r.interleaved=t,Mt(e,n)}function Pl(e,t,n){if(t=t.updateQueue,t!==null&&(t=t.shared,(n&4194240)!==0)){var r=t.lanes;r&=e.pendingLanes,n|=r,t.lanes=n,ho(e,n)}}function ga(e,t){var n=e.updateQueue,r=e.alternate;if(r!==null&&(r=r.updateQueue,n===r)){var l=null,o=null;if(n=n.firstBaseUpdate,n!==null){do{var i={eventTime:n.eventTime,lane:n.lane,tag:n.tag,payload:n.payload,callback:n.callback,next:null};o===null?l=o=i:o=o.next=i,n=n.next}while(n!==null);o===null?l=o=t:o=o.next=t}else l=o=t;n={baseState:r.baseState,firstBaseUpdate:l,lastBaseUpdate:o,shared:r.shared,effects:r.effects},e.updateQueue=n;return}e=n.lastBaseUpdate,e===null?n.firstBaseUpdate=t:e.next=t,n.lastBaseUpdate=t}function Tl(e,t,n,r){var l=e.updateQueue;tn=!1;var o=l.firstBaseUpdate,i=l.lastBaseUpdate,s=l.shared.pending;if(s!==null){l.shared.pending=null;var c=s,g=c.next;c.next=null,i===null?o=g:i.next=g,i=c;var S=e.alternate;S!==null&&(S=S.updateQueue,s=S.lastBaseUpdate,s!==i&&(s===null?S.firstBaseUpdate=g:s.next=g,S.lastBaseUpdate=c))}if(o!==null){var $=l.baseState;i=0,S=g=c=null,s=o;do{var k=s.lane,z=s.eventTime;if((r&k)===k){S!==null&&(S=S.next={eventTime:z,lane:0,tag:s.tag,payload:s.payload,callback:s.callback,next:null});e:{var O=e,F=s;switch(k=t,z=n,F.tag){case 1:if(O=F.payload,typeof O=="function"){$=O.call(z,$,k);break e}$=O;break e;case 3:O.flags=O.flags&-65537|128;case 0:if(O=F.payload,k=typeof O=="function"?O.call(z,$,k):O,k==null)break e;$=w({},$,k);break e;case 2:tn=!0}}s.callback!==null&&s.lane!==0&&(e.flags|=64,k=l.effects,k===null?l.effects=[s]:k.push(s))}else z={eventTime:z,lane:k,tag:s.tag,payload:s.payload,callback:s.callback,next:null},S===null?(g=S=z,c=$):S=S.next=z,i|=k;if(s=s.next,s===null){if(s=l.shared.pending,s===null)break;k=s,s=k.next,k.next=null,l.lastBaseUpdate=k,l.shared.pending=null}}while(!0);if(S===null&&(c=$),l.baseState=c,l.firstBaseUpdate=g,l.lastBaseUpdate=S,t=l.shared.interleaved,t!==null){l=t;do i|=l.lane,l=l.next;while(l!==t)}else o===null&&(l.shared.lanes=0);Cn|=i,e.lanes=i,e.memoizedState=$}}function ya(e,t,n){if(e=t.effects,t.effects=null,e!==null)for(t=0;t<e.length;t++){var r=e[t],l=r.callback;if(l!==null){if(r.callback=null,r=n,typeof l!="function")throw Error(m(191,l));l.call(r)}}}var Rr={},Tt=Jt(Rr),Lr=Jt(Rr),Ir=Jt(Rr);function Nn(e){if(e===Rr)throw Error(m(174));return e}function ri(e,t){switch(ce(Ir,t),ce(Lr,e),ce(Tt,Rr),e=t.nodeType,e){case 9:case 11:t=(t=t.documentElement)?t.namespaceURI:sr(null,"");break;default:e=e===8?t.parentNode:t,t=e.namespaceURI||null,e=e.tagName,t=sr(t,e)}de(Tt),ce(Tt,t)}function Zn(){de(Tt),de(Lr),de(Ir)}function _a(e){Nn(Ir.current);var t=Nn(Tt.current),n=sr(t,e.type);t!==n&&(ce(Lr,e),ce(Tt,n))}function li(e){Lr.current===e&&(de(Tt),de(Lr))}var ge=Jt(0);function zl(e){for(var t=e;t!==null;){if(t.tag===13){var n=t.memoizedState;if(n!==null&&(n=n.dehydrated,n===null||n.data==="$?"||n.data==="$!"))return t}else if(t.tag===19&&t.memoizedProps.revealOrder!==void 0){if((t.flags&128)!==0)return t}else if(t.child!==null){t.child.return=t,t=t.child;continue}if(t===e)break;for(;t.sibling===null;){if(t.return===null||t.return===e)return null;t=t.return}t.sibling.return=t.return,t=t.sibling}return null}var oi=[];function ii(){for(var e=0;e<oi.length;e++)oi[e]._workInProgressVersionPrimary=null;oi.length=0}var Rl=we.ReactCurrentDispatcher,si=we.ReactCurrentBatchConfig,En=0,ye=null,xe=null,Re=null,Ll=!1,Dr=!1,Or=0,pf=0;function Qe(){throw Error(m(321))}function ai(e,t){if(t===null)return!1;for(var n=0;n<t.length&&n<e.length;n++)if(!wt(e[n],t[n]))return!1;return!0}function ui(e,t,n,r,l,o){if(En=o,ye=t,t.memoizedState=null,t.updateQueue=null,t.lanes=0,Rl.current=e===null||e.memoizedState===null?gf:yf,e=n(r,l),Dr){o=0;do{if(Dr=!1,Or=0,25<=o)throw Error(m(301));o+=1,Re=xe=null,t.updateQueue=null,Rl.current=_f,e=n(r,l)}while(Dr)}if(Rl.current=Ol,t=xe!==null&&xe.next!==null,En=0,Re=xe=ye=null,Ll=!1,t)throw Error(m(300));return e}function ci(){var e=Or!==0;return Or=0,e}function zt(){var e={memoizedState:null,baseState:null,baseQueue:null,queue:null,next:null};return Re===null?ye.memoizedState=Re=e:Re=Re.next=e,Re}function pt(){if(xe===null){var e=ye.alternate;e=e!==null?e.memoizedState:null}else e=xe.next;var t=Re===null?ye.memoizedState:Re.next;if(t!==null)Re=t,xe=e;else{if(e===null)throw Error(m(310));xe=e,e={memoizedState:xe.memoizedState,baseState:xe.baseState,baseQueue:xe.baseQueue,queue:xe.queue,next:null},Re===null?ye.memoizedState=Re=e:Re=Re.next=e}return Re}function Fr(e,t){return typeof t=="function"?t(e):t}function fi(e){var t=pt(),n=t.queue;if(n===null)throw Error(m(311));n.lastRenderedReducer=e;var r=xe,l=r.baseQueue,o=n.pending;if(o!==null){if(l!==null){var i=l.next;l.next=o.next,o.next=i}r.baseQueue=l=o,n.pending=null}if(l!==null){o=l.next,r=r.baseState;var s=i=null,c=null,g=o;do{var S=g.lane;if((En&S)===S)c!==null&&(c=c.next={lane:0,action:g.action,hasEagerState:g.hasEagerState,eagerState:g.eagerState,next:null}),r=g.hasEagerState?g.eagerState:e(r,g.action);else{var $={lane:S,action:g.action,hasEagerState:g.hasEagerState,eagerState:g.eagerState,next:null};c===null?(s=c=$,i=r):c=c.next=$,ye.lanes|=S,Cn|=S}g=g.next}while(g!==null&&g!==o);c===null?i=r:c.next=s,wt(r,t.memoizedState)||(Ze=!0),t.memoizedState=r,t.baseState=i,t.baseQueue=c,n.lastRenderedState=r}if(e=n.interleaved,e!==null){l=e;do o=l.lane,ye.lanes|=o,Cn|=o,l=l.next;while(l!==e)}else l===null&&(n.lanes=0);return[t.memoizedState,n.dispatch]}function di(e){var t=pt(),n=t.queue;if(n===null)throw Error(m(311));n.lastRenderedReducer=e;var r=n.dispatch,l=n.pending,o=t.memoizedState;if(l!==null){n.pending=null;var i=l=l.next;do o=e(o,i.action),i=i.next;while(i!==l);wt(o,t.memoizedState)||(Ze=!0),t.memoizedState=o,t.baseQueue===null&&(t.baseState=o),n.lastRenderedState=o}return[o,r]}function wa(){}function ka(e,t){var n=ye,r=pt(),l=t(),o=!wt(r.memoizedState,l);if(o&&(r.memoizedState=l,Ze=!0),r=r.queue,pi(Na.bind(null,n,r,e),[e]),r.getSnapshot!==t||o||Re!==null&&Re.memoizedState.tag&1){if(n.flags|=2048,Mr(9,$a.bind(null,n,r,l,t),void 0,null),Le===null)throw Error(m(349));(En&30)!==0||Sa(n,t,l)}return l}function Sa(e,t,n){e.flags|=16384,e={getSnapshot:t,value:n},t=ye.updateQueue,t===null?(t={lastEffect:null,stores:null},ye.updateQueue=t,t.stores=[e]):(n=t.stores,n===null?t.stores=[e]:n.push(e))}function $a(e,t,n,r){t.value=n,t.getSnapshot=r,Ea(t)&&Ca(e)}function Na(e,t,n){return n(function(){Ea(t)&&Ca(e)})}function Ea(e){var t=e.getSnapshot;e=e.value;try{var n=t();return!wt(e,n)}catch{return!0}}function Ca(e){var t=Mt(e,1);t!==null&&Et(t,e,1,-1)}function xa(e){var t=zt();return typeof e=="function"&&(e=e()),t.memoizedState=t.baseState=e,e={pending:null,interleaved:null,lanes:0,dispatch:null,lastRenderedReducer:Fr,lastRenderedState:e},t.queue=e,e=e.dispatch=hf.bind(null,ye,e),[t.memoizedState,e]}function Mr(e,t,n,r){return e={tag:e,create:t,destroy:n,deps:r,next:null},t=ye.updateQueue,t===null?(t={lastEffect:null,stores:null},ye.updateQueue=t,t.lastEffect=e.next=e):(n=t.lastEffect,n===null?t.lastEffect=e.next=e:(r=n.next,n.next=e,e.next=r,t.lastEffect=e)),e}function Pa(){return pt().memoizedState}function Il(e,t,n,r){var l=zt();ye.flags|=e,l.memoizedState=Mr(1|t,n,void 0,r===void 0?null:r)}function Dl(e,t,n,r){var l=pt();r=r===void 0?null:r;var o=void 0;if(xe!==null){var i=xe.memoizedState;if(o=i.destroy,r!==null&&ai(r,i.deps)){l.memoizedState=Mr(t,n,o,r);return}}ye.flags|=e,l.memoizedState=Mr(1|t,n,o,r)}function Ta(e,t){return Il(8390656,8,e,t)}function pi(e,t){return Dl(2048,8,e,t)}function za(e,t){return Dl(4,2,e,t)}function Ra(e,t){return Dl(4,4,e,t)}function La(e,t){if(typeof t=="function")return e=e(),t(e),function(){t(null)};if(t!=null)return e=e(),t.current=e,function(){t.current=null}}function Ia(e,t,n){return n=n!=null?n.concat([e]):null,Dl(4,4,La.bind(null,t,e),n)}function mi(){}function Da(e,t){var n=pt();t=t===void 0?null:t;var r=n.memoizedState;return r!==null&&t!==null&&ai(t,r[1])?r[0]:(n.memoizedState=[e,t],e)}function Oa(e,t){var n=pt();t=t===void 0?null:t;var r=n.memoizedState;return r!==null&&t!==null&&ai(t,r[1])?r[0]:(e=e(),n.memoizedState=[e,t],e)}function Fa(e,t,n){return(En&21)===0?(e.baseState&&(e.baseState=!1,Ze=!0),e.memoizedState=n):(wt(n,t)||(n=ds(),ye.lanes|=n,Cn|=n,e.baseState=!0),t)}function mf(e,t){var n=ie;ie=n!==0&&4>n?n:4,e(!0);var r=si.transition;si.transition={};try{e(!1),t()}finally{ie=n,si.transition=r}}function Ma(){return pt().memoizedState}function vf(e,t,n){var r=sn(e);if(n={lane:r,action:n,hasEagerState:!1,eagerState:null,next:null},ja(e))Aa(t,n);else if(n=va(e,t,n,r),n!==null){var l=Ge();Et(n,e,r,l),Ua(n,t,r)}}function hf(e,t,n){var r=sn(e),l={lane:r,action:n,hasEagerState:!1,eagerState:null,next:null};if(ja(e))Aa(t,l);else{var o=e.alternate;if(e.lanes===0&&(o===null||o.lanes===0)&&(o=t.lastRenderedReducer,o!==null))try{var i=t.lastRenderedState,s=o(i,n);if(l.hasEagerState=!0,l.eagerState=s,wt(s,i)){var c=t.interleaved;c===null?(l.next=l,ti(t)):(l.next=c.next,c.next=l),t.interleaved=l;return}}catch{}n=va(e,t,l,r),n!==null&&(l=Ge(),Et(n,e,r,l),Ua(n,t,r))}}function ja(e){var t=e.alternate;return e===ye||t!==null&&t===ye}function Aa(e,t){Dr=Ll=!0;var n=e.pending;n===null?t.next=t:(t.next=n.next,n.next=t),e.pending=t}function Ua(e,t,n){if((n&4194240)!==0){var r=t.lanes;r&=e.pendingLanes,n|=r,t.lanes=n,ho(e,n)}}var Ol={readContext:dt,useCallback:Qe,useContext:Qe,useEffect:Qe,useImperativeHandle:Qe,useInsertionEffect:Qe,useLayoutEffect:Qe,useMemo:Qe,useReducer:Qe,useRef:Qe,useState:Qe,useDebugValue:Qe,useDeferredValue:Qe,useTransition:Qe,useMutableSource:Qe,useSyncExternalStore:Qe,useId:Qe,unstable_isNewReconciler:!1},gf={readContext:dt,useCallback:function(e,t){return zt().memoizedState=[e,t===void 0?null:t],e},useContext:dt,useEffect:Ta,useImperativeHandle:function(e,t,n){return n=n!=null?n.concat([e]):null,Il(4194308,4,La.bind(null,t,e),n)},useLayoutEffect:function(e,t){return Il(4194308,4,e,t)},useInsertionEffect:function(e,t){return Il(4,2,e,t)},useMemo:function(e,t){var n=zt();return t=t===void 0?null:t,e=e(),n.memoizedState=[e,t],e},useReducer:function(e,t,n){var r=zt();return t=n!==void 0?n(t):t,r.memoizedState=r.baseState=t,e={pending:null,interleaved:null,lanes:0,dispatch:null,lastRenderedReducer:e,lastRenderedState:t},r.queue=e,e=e.dispatch=vf.bind(null,ye,e),[r.memoizedState,e]},useRef:function(e){var t=zt();return e={current:e},t.memoizedState=e},useState:xa,useDebugValue:mi,useDeferredValue:function(e){return zt().memoizedState=e},useTransition:function(){var e=xa(!1),t=e[0];return e=mf.bind(null,e[1]),zt().memoizedState=e,[t,e]},useMutableSource:function(){},useSyncExternalStore:function(e,t,n){var r=ye,l=zt();if(ve){if(n===void 0)throw Error(m(407));n=n()}else{if(n=t(),Le===null)throw Error(m(349));(En&30)!==0||Sa(r,t,n)}l.memoizedState=n;var o={value:n,getSnapshot:t};return l.queue=o,Ta(Na.bind(null,r,o,e),[e]),r.flags|=2048,Mr(9,$a.bind(null,r,o,n,t),void 0,null),n},useId:function(){var e=zt(),t=Le.identifierPrefix;if(ve){var n=Ft,r=Ot;n=(r&~(1<<32-_t(r)-1)).toString(32)+n,t=":"+t+"R"+n,n=Or++,0<n&&(t+="H"+n.toString(32)),t+=":"}else n=pf++,t=":"+t+"r"+n.toString(32)+":";return e.memoizedState=t},unstable_isNewReconciler:!1},yf={readContext:dt,useCallback:Da,useContext:dt,useEffect:pi,useImperativeHandle:Ia,useInsertionEffect:za,useLayoutEffect:Ra,useMemo:Oa,useReducer:fi,useRef:Pa,useState:function(){return fi(Fr)},useDebugValue:mi,useDeferredValue:function(e){var t=pt();return Fa(t,xe.memoizedState,e)},useTransition:function(){var e=fi(Fr)[0],t=pt().memoizedState;return[e,t]},useMutableSource:wa,useSyncExternalStore:ka,useId:Ma,unstable_isNewReconciler:!1},_f={readContext:dt,useCallback:Da,useContext:dt,useEffect:pi,useImperativeHandle:Ia,useInsertionEffect:za,useLayoutEffect:Ra,useMemo:Oa,useReducer:di,useRef:Pa,useState:function(){return di(Fr)},useDebugValue:mi,useDeferredValue:function(e){var t=pt();return xe===null?t.memoizedState=e:Fa(t,xe.memoizedState,e)},useTransition:function(){var e=di(Fr)[0],t=pt().memoizedState;return[e,t]},useMutableSource:wa,useSyncExternalStore:ka,useId:Ma,unstable_isNewReconciler:!1};function St(e,t){if(e&&e.defaultProps){t=w({},t),e=e.defaultProps;for(var n in e)t[n]===void 0&&(t[n]=e[n]);return t}return t}function vi(e,t,n,r){t=e.memoizedState,n=n(r,t),n=n==null?t:w({},t,n),e.memoizedState=n,e.lanes===0&&(e.updateQueue.baseState=n)}var Fl={isMounted:function(e){return(e=e._reactInternals)?yn(e)===e:!1},enqueueSetState:function(e,t,n){e=e._reactInternals;var r=Ge(),l=sn(e),o=jt(r,l);o.payload=t,n!=null&&(o.callback=n),t=nn(e,o,l),t!==null&&(Et(t,e,l,r),Pl(t,e,l))},enqueueReplaceState:function(e,t,n){e=e._reactInternals;var r=Ge(),l=sn(e),o=jt(r,l);o.tag=1,o.payload=t,n!=null&&(o.callback=n),t=nn(e,o,l),t!==null&&(Et(t,e,l,r),Pl(t,e,l))},enqueueForceUpdate:function(e,t){e=e._reactInternals;var n=Ge(),r=sn(e),l=jt(n,r);l.tag=2,t!=null&&(l.callback=t),t=nn(e,l,r),t!==null&&(Et(t,e,r,n),Pl(t,e,r))}};function Ba(e,t,n,r,l,o,i){return e=e.stateNode,typeof e.shouldComponentUpdate=="function"?e.shouldComponentUpdate(r,o,i):t.prototype&&t.prototype.isPureReactComponent?!$r(n,r)||!$r(l,o):!0}function Wa(e,t,n){var r=!1,l=bt,o=t.contextType;return typeof o=="object"&&o!==null?o=dt(o):(l=qe(t)?wn:We.current,r=t.contextTypes,o=(r=r!=null)?Vn(e,l):bt),t=new t(n,o),e.memoizedState=t.state!==null&&t.state!==void 0?t.state:null,t.updater=Fl,e.stateNode=t,t._reactInternals=e,r&&(e=e.stateNode,e.__reactInternalMemoizedUnmaskedChildContext=l,e.__reactInternalMemoizedMaskedChildContext=o),t}function Qa(e,t,n,r){e=t.state,typeof t.componentWillReceiveProps=="function"&&t.componentWillReceiveProps(n,r),typeof t.UNSAFE_componentWillReceiveProps=="function"&&t.UNSAFE_componentWillReceiveProps(n,r),t.state!==e&&Fl.enqueueReplaceState(t,t.state,null)}function hi(e,t,n,r){var l=e.stateNode;l.props=n,l.state=e.memoizedState,l.refs={},ni(e);var o=t.contextType;typeof o=="object"&&o!==null?l.context=dt(o):(o=qe(t)?wn:We.current,l.context=Vn(e,o)),l.state=e.memoizedState,o=t.getDerivedStateFromProps,typeof o=="function"&&(vi(e,t,o,n),l.state=e.memoizedState),typeof t.getDerivedStateFromProps=="function"||typeof l.getSnapshotBeforeUpdate=="function"||typeof l.UNSAFE_componentWillMount!="function"&&typeof l.componentWillMount!="function"||(t=l.state,typeof l.componentWillMount=="function"&&l.componentWillMount(),typeof l.UNSAFE_componentWillMount=="function"&&l.UNSAFE_componentWillMount(),t!==l.state&&Fl.enqueueReplaceState(l,l.state,null),Tl(e,n,l,r),l.state=e.memoizedState),typeof l.componentDidMount=="function"&&(e.flags|=4194308)}function Jn(e,t){try{var n="",r=t;do n+=R(r),r=r.return;while(r);var l=n}catch(o){l=`
Error generating stack: `+o.message+`
`+o.stack}return{value:e,source:t,stack:l,digest:null}}function gi(e,t,n){return{value:e,source:null,stack:n??null,digest:t??null}}function yi(e,t){try{console.error(t.value)}catch(n){setTimeout(function(){throw n})}}var wf=typeof WeakMap=="function"?WeakMap:Map;function Va(e,t,n){n=jt(-1,n),n.tag=3,n.payload={element:null};var r=t.value;return n.callback=function(){Ql||(Ql=!0,Ii=r),yi(e,t)},n}function Ha(e,t,n){n=jt(-1,n),n.tag=3;var r=e.type.getDerivedStateFromError;if(typeof r=="function"){var l=t.value;n.payload=function(){return r(l)},n.callback=function(){yi(e,t)}}var o=e.stateNode;return o!==null&&typeof o.componentDidCatch=="function"&&(n.callback=function(){yi(e,t),typeof r!="function"&&(ln===null?ln=new Set([this]):ln.add(this));var i=t.stack;this.componentDidCatch(t.value,{componentStack:i!==null?i:""})}),n}function Ka(e,t,n){var r=e.pingCache;if(r===null){r=e.pingCache=new wf;var l=new Set;r.set(t,l)}else l=r.get(t),l===void 0&&(l=new Set,r.set(t,l));l.has(n)||(l.add(n),e=Df.bind(null,e,t,n),t.then(e,e))}function Ya(e){do{var t;if((t=e.tag===13)&&(t=e.memoizedState,t=t!==null?t.dehydrated!==null:!0),t)return e;e=e.return}while(e!==null);return null}function Ga(e,t,n,r,l){return(e.mode&1)===0?(e===t?e.flags|=65536:(e.flags|=128,n.flags|=131072,n.flags&=-52805,n.tag===1&&(n.alternate===null?n.tag=17:(t=jt(-1,1),t.tag=2,nn(n,t,1))),n.lanes|=1),e):(e.flags|=65536,e.lanes=l,e)}var kf=we.ReactCurrentOwner,Ze=!1;function Ye(e,t,n,r){t.child=e===null?ma(t,null,n,r):Gn(t,e.child,n,r)}function Xa(e,t,n,r,l){n=n.render;var o=t.ref;return qn(t,l),r=ui(e,t,n,r,o,l),n=ci(),e!==null&&!Ze?(t.updateQueue=e.updateQueue,t.flags&=-2053,e.lanes&=~l,At(e,t,l)):(ve&&n&&Ko(t),t.flags|=1,Ye(e,t,r,l),t.child)}function qa(e,t,n,r,l){if(e===null){var o=n.type;return typeof o=="function"&&!Ui(o)&&o.defaultProps===void 0&&n.compare===null&&n.defaultProps===void 0?(t.tag=15,t.type=o,Za(e,t,o,r,l)):(e=Xl(n.type,null,r,t,t.mode,l),e.ref=t.ref,e.return=t,t.child=e)}if(o=e.child,(e.lanes&l)===0){var i=o.memoizedProps;if(n=n.compare,n=n!==null?n:$r,n(i,r)&&e.ref===t.ref)return At(e,t,l)}return t.flags|=1,e=un(o,r),e.ref=t.ref,e.return=t,t.child=e}function Za(e,t,n,r,l){if(e!==null){var o=e.memoizedProps;if($r(o,r)&&e.ref===t.ref)if(Ze=!1,t.pendingProps=r=o,(e.lanes&l)!==0)(e.flags&131072)!==0&&(Ze=!0);else return t.lanes=e.lanes,At(e,t,l)}return _i(e,t,n,r,l)}function Ja(e,t,n){var r=t.pendingProps,l=r.children,o=e!==null?e.memoizedState:null;if(r.mode==="hidden")if((t.mode&1)===0)t.memoizedState={baseLanes:0,cachePool:null,transitions:null},ce(er,st),st|=n;else{if((n&1073741824)===0)return e=o!==null?o.baseLanes|n:n,t.lanes=t.childLanes=1073741824,t.memoizedState={baseLanes:e,cachePool:null,transitions:null},t.updateQueue=null,ce(er,st),st|=e,null;t.memoizedState={baseLanes:0,cachePool:null,transitions:null},r=o!==null?o.baseLanes:n,ce(er,st),st|=r}else o!==null?(r=o.baseLanes|n,t.memoizedState=null):r=n,ce(er,st),st|=r;return Ye(e,t,l,n),t.child}function ba(e,t){var n=t.ref;(e===null&&n!==null||e!==null&&e.ref!==n)&&(t.flags|=512,t.flags|=2097152)}function _i(e,t,n,r,l){var o=qe(n)?wn:We.current;return o=Vn(t,o),qn(t,l),n=ui(e,t,n,r,o,l),r=ci(),e!==null&&!Ze?(t.updateQueue=e.updateQueue,t.flags&=-2053,e.lanes&=~l,At(e,t,l)):(ve&&r&&Ko(t),t.flags|=1,Ye(e,t,n,l),t.child)}function eu(e,t,n,r,l){if(qe(n)){var o=!0;wl(t)}else o=!1;if(qn(t,l),t.stateNode===null)jl(e,t),Wa(t,n,r),hi(t,n,r,l),r=!0;else if(e===null){var i=t.stateNode,s=t.memoizedProps;i.props=s;var c=i.context,g=n.contextType;typeof g=="object"&&g!==null?g=dt(g):(g=qe(n)?wn:We.current,g=Vn(t,g));var S=n.getDerivedStateFromProps,$=typeof S=="function"||typeof i.getSnapshotBeforeUpdate=="function";$||typeof i.UNSAFE_componentWillReceiveProps!="function"&&typeof i.componentWillReceiveProps!="function"||(s!==r||c!==g)&&Qa(t,i,r,g),tn=!1;var k=t.memoizedState;i.state=k,Tl(t,r,i,l),c=t.memoizedState,s!==r||k!==c||Xe.current||tn?(typeof S=="function"&&(vi(t,n,S,r),c=t.memoizedState),(s=tn||Ba(t,n,s,r,k,c,g))?($||typeof i.UNSAFE_componentWillMount!="function"&&typeof i.componentWillMount!="function"||(typeof i.componentWillMount=="function"&&i.componentWillMount(),typeof i.UNSAFE_componentWillMount=="function"&&i.UNSAFE_componentWillMount()),typeof i.componentDidMount=="function"&&(t.flags|=4194308)):(typeof i.componentDidMount=="function"&&(t.flags|=4194308),t.memoizedProps=r,t.memoizedState=c),i.props=r,i.state=c,i.context=g,r=s):(typeof i.componentDidMount=="function"&&(t.flags|=4194308),r=!1)}else{i=t.stateNode,ha(e,t),s=t.memoizedProps,g=t.type===t.elementType?s:St(t.type,s),i.props=g,$=t.pendingProps,k=i.context,c=n.contextType,typeof c=="object"&&c!==null?c=dt(c):(c=qe(n)?wn:We.current,c=Vn(t,c));var z=n.getDerivedStateFromProps;(S=typeof z=="function"||typeof i.getSnapshotBeforeUpdate=="function")||typeof i.UNSAFE_componentWillReceiveProps!="function"&&typeof i.componentWillReceiveProps!="function"||(s!==$||k!==c)&&Qa(t,i,r,c),tn=!1,k=t.memoizedState,i.state=k,Tl(t,r,i,l);var O=t.memoizedState;s!==$||k!==O||Xe.current||tn?(typeof z=="function"&&(vi(t,n,z,r),O=t.memoizedState),(g=tn||Ba(t,n,g,r,k,O,c)||!1)?(S||typeof i.UNSAFE_componentWillUpdate!="function"&&typeof i.componentWillUpdate!="function"||(typeof i.componentWillUpdate=="function"&&i.componentWillUpdate(r,O,c),typeof i.UNSAFE_componentWillUpdate=="function"&&i.UNSAFE_componentWillUpdate(r,O,c)),typeof i.componentDidUpdate=="function"&&(t.flags|=4),typeof i.getSnapshotBeforeUpdate=="function"&&(t.flags|=1024)):(typeof i.componentDidUpdate!="function"||s===e.memoizedProps&&k===e.memoizedState||(t.flags|=4),typeof i.getSnapshotBeforeUpdate!="function"||s===e.memoizedProps&&k===e.memoizedState||(t.flags|=1024),t.memoizedProps=r,t.memoizedState=O),i.props=r,i.state=O,i.context=c,r=g):(typeof i.componentDidUpdate!="function"||s===e.memoizedProps&&k===e.memoizedState||(t.flags|=4),typeof i.getSnapshotBeforeUpdate!="function"||s===e.memoizedProps&&k===e.memoizedState||(t.flags|=1024),r=!1)}return wi(e,t,n,r,o,l)}function wi(e,t,n,r,l,o){ba(e,t);var i=(t.flags&128)!==0;if(!r&&!i)return l&&oa(t,n,!1),At(e,t,o);r=t.stateNode,kf.current=t;var s=i&&typeof n.getDerivedStateFromError!="function"?null:r.render();return t.flags|=1,e!==null&&i?(t.child=Gn(t,e.child,null,o),t.child=Gn(t,null,s,o)):Ye(e,t,s,o),t.memoizedState=r.state,l&&oa(t,n,!0),t.child}function tu(e){var t=e.stateNode;t.pendingContext?ra(e,t.pendingContext,t.pendingContext!==t.context):t.context&&ra(e,t.context,!1),ri(e,t.containerInfo)}function nu(e,t,n,r,l){return Yn(),qo(l),t.flags|=256,Ye(e,t,n,r),t.child}var ki={dehydrated:null,treeContext:null,retryLane:0};function Si(e){return{baseLanes:e,cachePool:null,transitions:null}}function ru(e,t,n){var r=t.pendingProps,l=ge.current,o=!1,i=(t.flags&128)!==0,s;if((s=i)||(s=e!==null&&e.memoizedState===null?!1:(l&2)!==0),s?(o=!0,t.flags&=-129):(e===null||e.memoizedState!==null)&&(l|=1),ce(ge,l&1),e===null)return Xo(t),e=t.memoizedState,e!==null&&(e=e.dehydrated,e!==null)?((t.mode&1)===0?t.lanes=1:e.data==="$!"?t.lanes=8:t.lanes=1073741824,null):(i=r.children,e=r.fallback,o?(r=t.mode,o=t.child,i={mode:"hidden",children:i},(r&1)===0&&o!==null?(o.childLanes=0,o.pendingProps=i):o=ql(i,r,0,null),e=zn(e,r,n,null),o.return=t,e.return=t,o.sibling=e,t.child=o,t.child.memoizedState=Si(n),t.memoizedState=ki,e):$i(t,i));if(l=e.memoizedState,l!==null&&(s=l.dehydrated,s!==null))return Sf(e,t,i,r,s,l,n);if(o){o=r.fallback,i=t.mode,l=e.child,s=l.sibling;var c={mode:"hidden",children:r.children};return(i&1)===0&&t.child!==l?(r=t.child,r.childLanes=0,r.pendingProps=c,t.deletions=null):(r=un(l,c),r.subtreeFlags=l.subtreeFlags&14680064),s!==null?o=un(s,o):(o=zn(o,i,n,null),o.flags|=2),o.return=t,r.return=t,r.sibling=o,t.child=r,r=o,o=t.child,i=e.child.memoizedState,i=i===null?Si(n):{baseLanes:i.baseLanes|n,cachePool:null,transitions:i.transitions},o.memoizedState=i,o.childLanes=e.childLanes&~n,t.memoizedState=ki,r}return o=e.child,e=o.sibling,r=un(o,{mode:"visible",children:r.children}),(t.mode&1)===0&&(r.lanes=n),r.return=t,r.sibling=null,e!==null&&(n=t.deletions,n===null?(t.deletions=[e],t.flags|=16):n.push(e)),t.child=r,t.memoizedState=null,r}function $i(e,t){return t=ql({mode:"visible",children:t},e.mode,0,null),t.return=e,e.child=t}function Ml(e,t,n,r){return r!==null&&qo(r),Gn(t,e.child,null,n),e=$i(t,t.pendingProps.children),e.flags|=2,t.memoizedState=null,e}function Sf(e,t,n,r,l,o,i){if(n)return t.flags&256?(t.flags&=-257,r=gi(Error(m(422))),Ml(e,t,i,r)):t.memoizedState!==null?(t.child=e.child,t.flags|=128,null):(o=r.fallback,l=t.mode,r=ql({mode:"visible",children:r.children},l,0,null),o=zn(o,l,i,null),o.flags|=2,r.return=t,o.return=t,r.sibling=o,t.child=r,(t.mode&1)!==0&&Gn(t,e.child,null,i),t.child.memoizedState=Si(i),t.memoizedState=ki,o);if((t.mode&1)===0)return Ml(e,t,i,null);if(l.data==="$!"){if(r=l.nextSibling&&l.nextSibling.dataset,r)var s=r.dgst;return r=s,o=Error(m(419)),r=gi(o,r,void 0),Ml(e,t,i,r)}if(s=(i&e.childLanes)!==0,Ze||s){if(r=Le,r!==null){switch(i&-i){case 4:l=2;break;case 16:l=8;break;case 64:case 128:case 256:case 512:case 1024:case 2048:case 4096:case 8192:case 16384:case 32768:case 65536:case 131072:case 262144:case 524288:case 1048576:case 2097152:case 4194304:case 8388608:case 16777216:case 33554432:case 67108864:l=32;break;case 536870912:l=268435456;break;default:l=0}l=(l&(r.suspendedLanes|i))!==0?0:l,l!==0&&l!==o.retryLane&&(o.retryLane=l,Mt(e,l),Et(r,e,l,-1))}return Ai(),r=gi(Error(m(421))),Ml(e,t,i,r)}return l.data==="$?"?(t.flags|=128,t.child=e.child,t=Of.bind(null,e),l._reactRetry=t,null):(e=o.treeContext,it=Zt(l.nextSibling),ot=t,ve=!0,kt=null,e!==null&&(ct[ft++]=Ot,ct[ft++]=Ft,ct[ft++]=kn,Ot=e.id,Ft=e.overflow,kn=t),t=$i(t,r.children),t.flags|=4096,t)}function lu(e,t,n){e.lanes|=t;var r=e.alternate;r!==null&&(r.lanes|=t),ei(e.return,t,n)}function Ni(e,t,n,r,l){var o=e.memoizedState;o===null?e.memoizedState={isBackwards:t,rendering:null,renderingStartTime:0,last:r,tail:n,tailMode:l}:(o.isBackwards=t,o.rendering=null,o.renderingStartTime=0,o.last=r,o.tail=n,o.tailMode=l)}function ou(e,t,n){var r=t.pendingProps,l=r.revealOrder,o=r.tail;if(Ye(e,t,r.children,n),r=ge.current,(r&2)!==0)r=r&1|2,t.flags|=128;else{if(e!==null&&(e.flags&128)!==0)e:for(e=t.child;e!==null;){if(e.tag===13)e.memoizedState!==null&&lu(e,n,t);else if(e.tag===19)lu(e,n,t);else if(e.child!==null){e.child.return=e,e=e.child;continue}if(e===t)break e;for(;e.sibling===null;){if(e.return===null||e.return===t)break e;e=e.return}e.sibling.return=e.return,e=e.sibling}r&=1}if(ce(ge,r),(t.mode&1)===0)t.memoizedState=null;else switch(l){case"forwards":for(n=t.child,l=null;n!==null;)e=n.alternate,e!==null&&zl(e)===null&&(l=n),n=n.sibling;n=l,n===null?(l=t.child,t.child=null):(l=n.sibling,n.sibling=null),Ni(t,!1,l,n,o);break;case"backwards":for(n=null,l=t.child,t.child=null;l!==null;){if(e=l.alternate,e!==null&&zl(e)===null){t.child=l;break}e=l.sibling,l.sibling=n,n=l,l=e}Ni(t,!0,n,null,o);break;case"together":Ni(t,!1,null,null,void 0);break;default:t.memoizedState=null}return t.child}function jl(e,t){(t.mode&1)===0&&e!==null&&(e.alternate=null,t.alternate=null,t.flags|=2)}function At(e,t,n){if(e!==null&&(t.dependencies=e.dependencies),Cn|=t.lanes,(n&t.childLanes)===0)return null;if(e!==null&&t.child!==e.child)throw Error(m(153));if(t.child!==null){for(e=t.child,n=un(e,e.pendingProps),t.child=n,n.return=t;e.sibling!==null;)e=e.sibling,n=n.sibling=un(e,e.pendingProps),n.return=t;n.sibling=null}return t.child}function $f(e,t,n){switch(t.tag){case 3:tu(t),Yn();break;case 5:_a(t);break;case 1:qe(t.type)&&wl(t);break;case 4:ri(t,t.stateNode.containerInfo);break;case 10:var r=t.type._context,l=t.memoizedProps.value;ce(Cl,r._currentValue),r._currentValue=l;break;case 13:if(r=t.memoizedState,r!==null)return r.dehydrated!==null?(ce(ge,ge.current&1),t.flags|=128,null):(n&t.child.childLanes)!==0?ru(e,t,n):(ce(ge,ge.current&1),e=At(e,t,n),e!==null?e.sibling:null);ce(ge,ge.current&1);break;case 19:if(r=(n&t.childLanes)!==0,(e.flags&128)!==0){if(r)return ou(e,t,n);t.flags|=128}if(l=t.memoizedState,l!==null&&(l.rendering=null,l.tail=null,l.lastEffect=null),ce(ge,ge.current),r)break;return null;case 22:case 23:return t.lanes=0,Ja(e,t,n)}return At(e,t,n)}var iu,Ei,su,au;iu=function(e,t){for(var n=t.child;n!==null;){if(n.tag===5||n.tag===6)e.appendChild(n.stateNode);else if(n.tag!==4&&n.child!==null){n.child.return=n,n=n.child;continue}if(n===t)break;for(;n.sibling===null;){if(n.return===null||n.return===t)return;n=n.return}n.sibling.return=n.return,n=n.sibling}},Ei=function(){},su=function(e,t,n,r){var l=e.memoizedProps;if(l!==r){e=t.stateNode,Nn(Tt.current);var o=null;switch(n){case"input":l=Wt(e,l),r=Wt(e,r),o=[];break;case"select":l=w({},l,{value:void 0}),r=w({},r,{value:void 0}),o=[];break;case"textarea":l=Lt(e,l),r=Lt(e,r),o=[];break;default:typeof l.onClick!="function"&&typeof r.onClick=="function"&&(e.onclick=gl)}me(n,r);var i;n=null;for(g in l)if(!r.hasOwnProperty(g)&&l.hasOwnProperty(g)&&l[g]!=null)if(g==="style"){var s=l[g];for(i in s)s.hasOwnProperty(i)&&(n||(n={}),n[i]="")}else g!=="dangerouslySetInnerHTML"&&g!=="children"&&g!=="suppressContentEditableWarning"&&g!=="suppressHydrationWarning"&&g!=="autoFocus"&&(C.hasOwnProperty(g)?o||(o=[]):(o=o||[]).push(g,null));for(g in r){var c=r[g];if(s=l?.[g],r.hasOwnProperty(g)&&c!==s&&(c!=null||s!=null))if(g==="style")if(s){for(i in s)!s.hasOwnProperty(i)||c&&c.hasOwnProperty(i)||(n||(n={}),n[i]="");for(i in c)c.hasOwnProperty(i)&&s[i]!==c[i]&&(n||(n={}),n[i]=c[i])}else n||(o||(o=[]),o.push(g,n)),n=c;else g==="dangerouslySetInnerHTML"?(c=c?c.__html:void 0,s=s?s.__html:void 0,c!=null&&s!==c&&(o=o||[]).push(g,c)):g==="children"?typeof c!="string"&&typeof c!="number"||(o=o||[]).push(g,""+c):g!=="suppressContentEditableWarning"&&g!=="suppressHydrationWarning"&&(C.hasOwnProperty(g)?(c!=null&&g==="onScroll"&&fe("scroll",e),o||s===c||(o=[])):(o=o||[]).push(g,c))}n&&(o=o||[]).push("style",n);var g=o;(t.updateQueue=g)&&(t.flags|=4)}},au=function(e,t,n,r){n!==r&&(t.flags|=4)};function jr(e,t){if(!ve)switch(e.tailMode){case"hidden":t=e.tail;for(var n=null;t!==null;)t.alternate!==null&&(n=t),t=t.sibling;n===null?e.tail=null:n.sibling=null;break;case"collapsed":n=e.tail;for(var r=null;n!==null;)n.alternate!==null&&(r=n),n=n.sibling;r===null?t||e.tail===null?e.tail=null:e.tail.sibling=null:r.sibling=null}}function Ve(e){var t=e.alternate!==null&&e.alternate.child===e.child,n=0,r=0;if(t)for(var l=e.child;l!==null;)n|=l.lanes|l.childLanes,r|=l.subtreeFlags&14680064,r|=l.flags&14680064,l.return=e,l=l.sibling;else for(l=e.child;l!==null;)n|=l.lanes|l.childLanes,r|=l.subtreeFlags,r|=l.flags,l.return=e,l=l.sibling;return e.subtreeFlags|=r,e.childLanes=n,t}function Nf(e,t,n){var r=t.pendingProps;switch(Yo(t),t.tag){case 2:case 16:case 15:case 0:case 11:case 7:case 8:case 12:case 9:case 14:return Ve(t),null;case 1:return qe(t.type)&&_l(),Ve(t),null;case 3:return r=t.stateNode,Zn(),de(Xe),de(We),ii(),r.pendingContext&&(r.context=r.pendingContext,r.pendingContext=null),(e===null||e.child===null)&&(Nl(t)?t.flags|=4:e===null||e.memoizedState.isDehydrated&&(t.flags&256)===0||(t.flags|=1024,kt!==null&&(Fi(kt),kt=null))),Ei(e,t),Ve(t),null;case 5:li(t);var l=Nn(Ir.current);if(n=t.type,e!==null&&t.stateNode!=null)su(e,t,n,r,l),e.ref!==t.ref&&(t.flags|=512,t.flags|=2097152);else{if(!r){if(t.stateNode===null)throw Error(m(166));return Ve(t),null}if(e=Nn(Tt.current),Nl(t)){r=t.stateNode,n=t.type;var o=t.memoizedProps;switch(r[Pt]=t,r[Pr]=o,e=(t.mode&1)!==0,n){case"dialog":fe("cancel",r),fe("close",r);break;case"iframe":case"object":case"embed":fe("load",r);break;case"video":case"audio":for(l=0;l<Er.length;l++)fe(Er[l],r);break;case"source":fe("error",r);break;case"img":case"image":case"link":fe("error",r),fe("load",r);break;case"details":fe("toggle",r);break;case"input":Vr(r,o),fe("invalid",r);break;case"select":r._wrapperState={wasMultiple:!!o.multiple},fe("invalid",r);break;case"textarea":Hr(r,o),fe("invalid",r)}me(n,o),l=null;for(var i in o)if(o.hasOwnProperty(i)){var s=o[i];i==="children"?typeof s=="string"?r.textContent!==s&&(o.suppressHydrationWarning!==!0&&hl(r.textContent,s,e),l=["children",s]):typeof s=="number"&&r.textContent!==""+s&&(o.suppressHydrationWarning!==!0&&hl(r.textContent,s,e),l=["children",""+s]):C.hasOwnProperty(i)&&s!=null&&i==="onScroll"&&fe("scroll",r)}switch(n){case"input":oe(r),Rn(r,o,!0);break;case"textarea":oe(r),Yr(r);break;case"select":case"option":break;default:typeof o.onClick=="function"&&(r.onclick=gl)}r=l,t.updateQueue=r,r!==null&&(t.flags|=4)}else{i=l.nodeType===9?l:l.ownerDocument,e==="http://www.w3.org/1999/xhtml"&&(e=Gr(n)),e==="http://www.w3.org/1999/xhtml"?n==="script"?(e=i.createElement("div"),e.innerHTML="<script><\/script>",e=e.removeChild(e.firstChild)):typeof r.is=="string"?e=i.createElement(n,{is:r.is}):(e=i.createElement(n),n==="select"&&(i=e,r.multiple?i.multiple=!0:r.size&&(i.size=r.size))):e=i.createElementNS(e,n),e[Pt]=t,e[Pr]=r,iu(e,t,!1,!1),t.stateNode=e;e:{switch(i=Vt(n,r),n){case"dialog":fe("cancel",e),fe("close",e),l=r;break;case"iframe":case"object":case"embed":fe("load",e),l=r;break;case"video":case"audio":for(l=0;l<Er.length;l++)fe(Er[l],e);l=r;break;case"source":fe("error",e),l=r;break;case"img":case"image":case"link":fe("error",e),fe("load",e),l=r;break;case"details":fe("toggle",e),l=r;break;case"input":Vr(e,r),l=Wt(e,r),fe("invalid",e);break;case"option":l=r;break;case"select":e._wrapperState={wasMultiple:!!r.multiple},l=w({},r,{value:void 0}),fe("invalid",e);break;case"textarea":Hr(e,r),l=Lt(e,r),fe("invalid",e);break;default:l=r}me(n,l),s=l;for(o in s)if(s.hasOwnProperty(o)){var c=s[o];o==="style"?L(e,c):o==="dangerouslySetInnerHTML"?(c=c?c.__html:void 0,c!=null&&Xr(e,c)):o==="children"?typeof c=="string"?(n!=="textarea"||c!=="")&&hn(e,c):typeof c=="number"&&hn(e,""+c):o!=="suppressContentEditableWarning"&&o!=="suppressHydrationWarning"&&o!=="autoFocus"&&(C.hasOwnProperty(o)?c!=null&&o==="onScroll"&&fe("scroll",e):c!=null&&nt(e,o,c,i))}switch(n){case"input":oe(e),Rn(e,r,!1);break;case"textarea":oe(e),Yr(e);break;case"option":r.value!=null&&e.setAttribute("value",""+b(r.value));break;case"select":e.multiple=!!r.multiple,o=r.value,o!=null?Qt(e,!!r.multiple,o,!1):r.defaultValue!=null&&Qt(e,!!r.multiple,r.defaultValue,!0);break;default:typeof l.onClick=="function"&&(e.onclick=gl)}switch(n){case"button":case"input":case"select":case"textarea":r=!!r.autoFocus;break e;case"img":r=!0;break e;default:r=!1}}r&&(t.flags|=4)}t.ref!==null&&(t.flags|=512,t.flags|=2097152)}return Ve(t),null;case 6:if(e&&t.stateNode!=null)au(e,t,e.memoizedProps,r);else{if(typeof r!="string"&&t.stateNode===null)throw Error(m(166));if(n=Nn(Ir.current),Nn(Tt.current),Nl(t)){if(r=t.stateNode,n=t.memoizedProps,r[Pt]=t,(o=r.nodeValue!==n)&&(e=ot,e!==null))switch(e.tag){case 3:hl(r.nodeValue,n,(e.mode&1)!==0);break;case 5:e.memoizedProps.suppressHydrationWarning!==!0&&hl(r.nodeValue,n,(e.mode&1)!==0)}o&&(t.flags|=4)}else r=(n.nodeType===9?n:n.ownerDocument).createTextNode(r),r[Pt]=t,t.stateNode=r}return Ve(t),null;case 13:if(de(ge),r=t.memoizedState,e===null||e.memoizedState!==null&&e.memoizedState.dehydrated!==null){if(ve&&it!==null&&(t.mode&1)!==0&&(t.flags&128)===0)fa(),Yn(),t.flags|=98560,o=!1;else if(o=Nl(t),r!==null&&r.dehydrated!==null){if(e===null){if(!o)throw Error(m(318));if(o=t.memoizedState,o=o!==null?o.dehydrated:null,!o)throw Error(m(317));o[Pt]=t}else Yn(),(t.flags&128)===0&&(t.memoizedState=null),t.flags|=4;Ve(t),o=!1}else kt!==null&&(Fi(kt),kt=null),o=!0;if(!o)return t.flags&65536?t:null}return(t.flags&128)!==0?(t.lanes=n,t):(r=r!==null,r!==(e!==null&&e.memoizedState!==null)&&r&&(t.child.flags|=8192,(t.mode&1)!==0&&(e===null||(ge.current&1)!==0?Pe===0&&(Pe=3):Ai())),t.updateQueue!==null&&(t.flags|=4),Ve(t),null);case 4:return Zn(),Ei(e,t),e===null&&Cr(t.stateNode.containerInfo),Ve(t),null;case 10:return bo(t.type._context),Ve(t),null;case 17:return qe(t.type)&&_l(),Ve(t),null;case 19:if(de(ge),o=t.memoizedState,o===null)return Ve(t),null;if(r=(t.flags&128)!==0,i=o.rendering,i===null)if(r)jr(o,!1);else{if(Pe!==0||e!==null&&(e.flags&128)!==0)for(e=t.child;e!==null;){if(i=zl(e),i!==null){for(t.flags|=128,jr(o,!1),r=i.updateQueue,r!==null&&(t.updateQueue=r,t.flags|=4),t.subtreeFlags=0,r=n,n=t.child;n!==null;)o=n,e=r,o.flags&=14680066,i=o.alternate,i===null?(o.childLanes=0,o.lanes=e,o.child=null,o.subtreeFlags=0,o.memoizedProps=null,o.memoizedState=null,o.updateQueue=null,o.dependencies=null,o.stateNode=null):(o.childLanes=i.childLanes,o.lanes=i.lanes,o.child=i.child,o.subtreeFlags=0,o.deletions=null,o.memoizedProps=i.memoizedProps,o.memoizedState=i.memoizedState,o.updateQueue=i.updateQueue,o.type=i.type,e=i.dependencies,o.dependencies=e===null?null:{lanes:e.lanes,firstContext:e.firstContext}),n=n.sibling;return ce(ge,ge.current&1|2),t.child}e=e.sibling}o.tail!==null&&ke()>tr&&(t.flags|=128,r=!0,jr(o,!1),t.lanes=4194304)}else{if(!r)if(e=zl(i),e!==null){if(t.flags|=128,r=!0,n=e.updateQueue,n!==null&&(t.updateQueue=n,t.flags|=4),jr(o,!0),o.tail===null&&o.tailMode==="hidden"&&!i.alternate&&!ve)return Ve(t),null}else 2*ke()-o.renderingStartTime>tr&&n!==1073741824&&(t.flags|=128,r=!0,jr(o,!1),t.lanes=4194304);o.isBackwards?(i.sibling=t.child,t.child=i):(n=o.last,n!==null?n.sibling=i:t.child=i,o.last=i)}return o.tail!==null?(t=o.tail,o.rendering=t,o.tail=t.sibling,o.renderingStartTime=ke(),t.sibling=null,n=ge.current,ce(ge,r?n&1|2:n&1),t):(Ve(t),null);case 22:case 23:return ji(),r=t.memoizedState!==null,e!==null&&e.memoizedState!==null!==r&&(t.flags|=8192),r&&(t.mode&1)!==0?(st&1073741824)!==0&&(Ve(t),t.subtreeFlags&6&&(t.flags|=8192)):Ve(t),null;case 24:return null;case 25:return null}throw Error(m(156,t.tag))}function Ef(e,t){switch(Yo(t),t.tag){case 1:return qe(t.type)&&_l(),e=t.flags,e&65536?(t.flags=e&-65537|128,t):null;case 3:return Zn(),de(Xe),de(We),ii(),e=t.flags,(e&65536)!==0&&(e&128)===0?(t.flags=e&-65537|128,t):null;case 5:return li(t),null;case 13:if(de(ge),e=t.memoizedState,e!==null&&e.dehydrated!==null){if(t.alternate===null)throw Error(m(340));Yn()}return e=t.flags,e&65536?(t.flags=e&-65537|128,t):null;case 19:return de(ge),null;case 4:return Zn(),null;case 10:return bo(t.type._context),null;case 22:case 23:return ji(),null;case 24:return null;default:return null}}var Al=!1,He=!1,Cf=typeof WeakSet=="function"?WeakSet:Set,I=null;function bn(e,t){var n=e.ref;if(n!==null)if(typeof n=="function")try{n(null)}catch(r){_e(e,t,r)}else n.current=null}function Ci(e,t,n){try{n()}catch(r){_e(e,t,r)}}var uu=!1;function xf(e,t){if(jo=ol,e=Bs(),zo(e)){if("selectionStart"in e)var n={start:e.selectionStart,end:e.selectionEnd};else e:{n=(n=e.ownerDocument)&&n.defaultView||window;var r=n.getSelection&&n.getSelection();if(r&&r.rangeCount!==0){n=r.anchorNode;var l=r.anchorOffset,o=r.focusNode;r=r.focusOffset;try{n.nodeType,o.nodeType}catch{n=null;break e}var i=0,s=-1,c=-1,g=0,S=0,$=e,k=null;t:for(;;){for(var z;$!==n||l!==0&&$.nodeType!==3||(s=i+l),$!==o||r!==0&&$.nodeType!==3||(c=i+r),$.nodeType===3&&(i+=$.nodeValue.length),(z=$.firstChild)!==null;)k=$,$=z;for(;;){if($===e)break t;if(k===n&&++g===l&&(s=i),k===o&&++S===r&&(c=i),(z=$.nextSibling)!==null)break;$=k,k=$.parentNode}$=z}n=s===-1||c===-1?null:{start:s,end:c}}else n=null}n=n||{start:0,end:0}}else n=null;for(Ao={focusedElem:e,selectionRange:n},ol=!1,I=t;I!==null;)if(t=I,e=t.child,(t.subtreeFlags&1028)!==0&&e!==null)e.return=t,I=e;else for(;I!==null;){t=I;try{var O=t.alternate;if((t.flags&1024)!==0)switch(t.tag){case 0:case 11:case 15:break;case 1:if(O!==null){var F=O.memoizedProps,Se=O.memoizedState,v=t.stateNode,p=v.getSnapshotBeforeUpdate(t.elementType===t.type?F:St(t.type,F),Se);v.__reactInternalSnapshotBeforeUpdate=p}break;case 3:var h=t.stateNode.containerInfo;h.nodeType===1?h.textContent="":h.nodeType===9&&h.documentElement&&h.removeChild(h.documentElement);break;case 5:case 6:case 4:case 17:break;default:throw Error(m(163))}}catch(N){_e(t,t.return,N)}if(e=t.sibling,e!==null){e.return=t.return,I=e;break}I=t.return}return O=uu,uu=!1,O}function Ar(e,t,n){var r=t.updateQueue;if(r=r!==null?r.lastEffect:null,r!==null){var l=r=r.next;do{if((l.tag&e)===e){var o=l.destroy;l.destroy=void 0,o!==void 0&&Ci(t,n,o)}l=l.next}while(l!==r)}}function Ul(e,t){if(t=t.updateQueue,t=t!==null?t.lastEffect:null,t!==null){var n=t=t.next;do{if((n.tag&e)===e){var r=n.create;n.destroy=r()}n=n.next}while(n!==t)}}function xi(e){var t=e.ref;if(t!==null){var n=e.stateNode;e.tag,e=n,typeof t=="function"?t(e):t.current=e}}function cu(e){var t=e.alternate;t!==null&&(e.alternate=null,cu(t)),e.child=null,e.deletions=null,e.sibling=null,e.tag===5&&(t=e.stateNode,t!==null&&(delete t[Pt],delete t[Pr],delete t[Qo],delete t[uf],delete t[cf])),e.stateNode=null,e.return=null,e.dependencies=null,e.memoizedProps=null,e.memoizedState=null,e.pendingProps=null,e.stateNode=null,e.updateQueue=null}function fu(e){return e.tag===5||e.tag===3||e.tag===4}function du(e){e:for(;;){for(;e.sibling===null;){if(e.return===null||fu(e.return))return null;e=e.return}for(e.sibling.return=e.return,e=e.sibling;e.tag!==5&&e.tag!==6&&e.tag!==18;){if(e.flags&2||e.child===null||e.tag===4)continue e;e.child.return=e,e=e.child}if(!(e.flags&2))return e.stateNode}}function Pi(e,t,n){var r=e.tag;if(r===5||r===6)e=e.stateNode,t?n.nodeType===8?n.parentNode.insertBefore(e,t):n.insertBefore(e,t):(n.nodeType===8?(t=n.parentNode,t.insertBefore(e,n)):(t=n,t.appendChild(e)),n=n._reactRootContainer,n!=null||t.onclick!==null||(t.onclick=gl));else if(r!==4&&(e=e.child,e!==null))for(Pi(e,t,n),e=e.sibling;e!==null;)Pi(e,t,n),e=e.sibling}function Ti(e,t,n){var r=e.tag;if(r===5||r===6)e=e.stateNode,t?n.insertBefore(e,t):n.appendChild(e);else if(r!==4&&(e=e.child,e!==null))for(Ti(e,t,n),e=e.sibling;e!==null;)Ti(e,t,n),e=e.sibling}var je=null,$t=!1;function rn(e,t,n){for(n=n.child;n!==null;)pu(e,t,n),n=n.sibling}function pu(e,t,n){if(xt&&typeof xt.onCommitFiberUnmount=="function")try{xt.onCommitFiberUnmount(br,n)}catch{}switch(n.tag){case 5:He||bn(n,t);case 6:var r=je,l=$t;je=null,rn(e,t,n),je=r,$t=l,je!==null&&($t?(e=je,n=n.stateNode,e.nodeType===8?e.parentNode.removeChild(n):e.removeChild(n)):je.removeChild(n.stateNode));break;case 18:je!==null&&($t?(e=je,n=n.stateNode,e.nodeType===8?Wo(e.parentNode,n):e.nodeType===1&&Wo(e,n),gr(e)):Wo(je,n.stateNode));break;case 4:r=je,l=$t,je=n.stateNode.containerInfo,$t=!0,rn(e,t,n),je=r,$t=l;break;case 0:case 11:case 14:case 15:if(!He&&(r=n.updateQueue,r!==null&&(r=r.lastEffect,r!==null))){l=r=r.next;do{var o=l,i=o.destroy;o=o.tag,i!==void 0&&((o&2)!==0||(o&4)!==0)&&Ci(n,t,i),l=l.next}while(l!==r)}rn(e,t,n);break;case 1:if(!He&&(bn(n,t),r=n.stateNode,typeof r.componentWillUnmount=="function"))try{r.props=n.memoizedProps,r.state=n.memoizedState,r.componentWillUnmount()}catch(s){_e(n,t,s)}rn(e,t,n);break;case 21:rn(e,t,n);break;case 22:n.mode&1?(He=(r=He)||n.memoizedState!==null,rn(e,t,n),He=r):rn(e,t,n);break;default:rn(e,t,n)}}function mu(e){var t=e.updateQueue;if(t!==null){e.updateQueue=null;var n=e.stateNode;n===null&&(n=e.stateNode=new Cf),t.forEach(function(r){var l=Ff.bind(null,e,r);n.has(r)||(n.add(r),r.then(l,l))})}}function Nt(e,t){var n=t.deletions;if(n!==null)for(var r=0;r<n.length;r++){var l=n[r];try{var o=e,i=t,s=i;e:for(;s!==null;){switch(s.tag){case 5:je=s.stateNode,$t=!1;break e;case 3:je=s.stateNode.containerInfo,$t=!0;break e;case 4:je=s.stateNode.containerInfo,$t=!0;break e}s=s.return}if(je===null)throw Error(m(160));pu(o,i,l),je=null,$t=!1;var c=l.alternate;c!==null&&(c.return=null),l.return=null}catch(g){_e(l,t,g)}}if(t.subtreeFlags&12854)for(t=t.child;t!==null;)vu(t,e),t=t.sibling}function vu(e,t){var n=e.alternate,r=e.flags;switch(e.tag){case 0:case 11:case 14:case 15:if(Nt(t,e),Rt(e),r&4){try{Ar(3,e,e.return),Ul(3,e)}catch(F){_e(e,e.return,F)}try{Ar(5,e,e.return)}catch(F){_e(e,e.return,F)}}break;case 1:Nt(t,e),Rt(e),r&512&&n!==null&&bn(n,n.return);break;case 5:if(Nt(t,e),Rt(e),r&512&&n!==null&&bn(n,n.return),e.flags&32){var l=e.stateNode;try{hn(l,"")}catch(F){_e(e,e.return,F)}}if(r&4&&(l=e.stateNode,l!=null)){var o=e.memoizedProps,i=n!==null?n.memoizedProps:o,s=e.type,c=e.updateQueue;if(e.updateQueue=null,c!==null)try{s==="input"&&o.type==="radio"&&o.name!=null&&or(l,o),Vt(s,i);var g=Vt(s,o);for(i=0;i<c.length;i+=2){var S=c[i],$=c[i+1];S==="style"?L(l,$):S==="dangerouslySetInnerHTML"?Xr(l,$):S==="children"?hn(l,$):nt(l,S,$,g)}switch(s){case"input":mn(l,o);break;case"textarea":Kr(l,o);break;case"select":var k=l._wrapperState.wasMultiple;l._wrapperState.wasMultiple=!!o.multiple;var z=o.value;z!=null?Qt(l,!!o.multiple,z,!1):k!==!!o.multiple&&(o.defaultValue!=null?Qt(l,!!o.multiple,o.defaultValue,!0):Qt(l,!!o.multiple,o.multiple?[]:"",!1))}l[Pr]=o}catch(F){_e(e,e.return,F)}}break;case 6:if(Nt(t,e),Rt(e),r&4){if(e.stateNode===null)throw Error(m(162));l=e.stateNode,o=e.memoizedProps;try{l.nodeValue=o}catch(F){_e(e,e.return,F)}}break;case 3:if(Nt(t,e),Rt(e),r&4&&n!==null&&n.memoizedState.isDehydrated)try{gr(t.containerInfo)}catch(F){_e(e,e.return,F)}break;case 4:Nt(t,e),Rt(e);break;case 13:Nt(t,e),Rt(e),l=e.child,l.flags&8192&&(o=l.memoizedState!==null,l.stateNode.isHidden=o,!o||l.alternate!==null&&l.alternate.memoizedState!==null||(Li=ke())),r&4&&mu(e);break;case 22:if(S=n!==null&&n.memoizedState!==null,e.mode&1?(He=(g=He)||S,Nt(t,e),He=g):Nt(t,e),Rt(e),r&8192){if(g=e.memoizedState!==null,(e.stateNode.isHidden=g)&&!S&&(e.mode&1)!==0)for(I=e,S=e.child;S!==null;){for($=I=S;I!==null;){switch(k=I,z=k.child,k.tag){case 0:case 11:case 14:case 15:Ar(4,k,k.return);break;case 1:bn(k,k.return);var O=k.stateNode;if(typeof O.componentWillUnmount=="function"){r=k,n=k.return;try{t=r,O.props=t.memoizedProps,O.state=t.memoizedState,O.componentWillUnmount()}catch(F){_e(r,n,F)}}break;case 5:bn(k,k.return);break;case 22:if(k.memoizedState!==null){yu($);continue}}z!==null?(z.return=k,I=z):yu($)}S=S.sibling}e:for(S=null,$=e;;){if($.tag===5){if(S===null){S=$;try{l=$.stateNode,g?(o=l.style,typeof o.setProperty=="function"?o.setProperty("display","none","important"):o.display="none"):(s=$.stateNode,c=$.memoizedProps.style,i=c!=null&&c.hasOwnProperty("display")?c.display:null,s.style.display=u("display",i))}catch(F){_e(e,e.return,F)}}}else if($.tag===6){if(S===null)try{$.stateNode.nodeValue=g?"":$.memoizedProps}catch(F){_e(e,e.return,F)}}else if(($.tag!==22&&$.tag!==23||$.memoizedState===null||$===e)&&$.child!==null){$.child.return=$,$=$.child;continue}if($===e)break e;for(;$.sibling===null;){if($.return===null||$.return===e)break e;S===$&&(S=null),$=$.return}S===$&&(S=null),$.sibling.return=$.return,$=$.sibling}}break;case 19:Nt(t,e),Rt(e),r&4&&mu(e);break;case 21:break;default:Nt(t,e),Rt(e)}}function Rt(e){var t=e.flags;if(t&2){try{e:{for(var n=e.return;n!==null;){if(fu(n)){var r=n;break e}n=n.return}throw Error(m(160))}switch(r.tag){case 5:var l=r.stateNode;r.flags&32&&(hn(l,""),r.flags&=-33);var o=du(e);Ti(e,o,l);break;case 3:case 4:var i=r.stateNode.containerInfo,s=du(e);Pi(e,s,i);break;default:throw Error(m(161))}}catch(c){_e(e,e.return,c)}e.flags&=-3}t&4096&&(e.flags&=-4097)}function Pf(e,t,n){I=e,hu(e)}function hu(e,t,n){for(var r=(e.mode&1)!==0;I!==null;){var l=I,o=l.child;if(l.tag===22&&r){var i=l.memoizedState!==null||Al;if(!i){var s=l.alternate,c=s!==null&&s.memoizedState!==null||He;s=Al;var g=He;if(Al=i,(He=c)&&!g)for(I=l;I!==null;)i=I,c=i.child,i.tag===22&&i.memoizedState!==null?_u(l):c!==null?(c.return=i,I=c):_u(l);for(;o!==null;)I=o,hu(o),o=o.sibling;I=l,Al=s,He=g}gu(e)}else(l.subtreeFlags&8772)!==0&&o!==null?(o.return=l,I=o):gu(e)}}function gu(e){for(;I!==null;){var t=I;if((t.flags&8772)!==0){var n=t.alternate;try{if((t.flags&8772)!==0)switch(t.tag){case 0:case 11:case 15:He||Ul(5,t);break;case 1:var r=t.stateNode;if(t.flags&4&&!He)if(n===null)r.componentDidMount();else{var l=t.elementType===t.type?n.memoizedProps:St(t.type,n.memoizedProps);r.componentDidUpdate(l,n.memoizedState,r.__reactInternalSnapshotBeforeUpdate)}var o=t.updateQueue;o!==null&&ya(t,o,r);break;case 3:var i=t.updateQueue;if(i!==null){if(n=null,t.child!==null)switch(t.child.tag){case 5:n=t.child.stateNode;break;case 1:n=t.child.stateNode}ya(t,i,n)}break;case 5:var s=t.stateNode;if(n===null&&t.flags&4){n=s;var c=t.memoizedProps;switch(t.type){case"button":case"input":case"select":case"textarea":c.autoFocus&&n.focus();break;case"img":c.src&&(n.src=c.src)}}break;case 6:break;case 4:break;case 12:break;case 13:if(t.memoizedState===null){var g=t.alternate;if(g!==null){var S=g.memoizedState;if(S!==null){var $=S.dehydrated;$!==null&&gr($)}}}break;case 19:case 17:case 21:case 22:case 23:case 25:break;default:throw Error(m(163))}He||t.flags&512&&xi(t)}catch(k){_e(t,t.return,k)}}if(t===e){I=null;break}if(n=t.sibling,n!==null){n.return=t.return,I=n;break}I=t.return}}function yu(e){for(;I!==null;){var t=I;if(t===e){I=null;break}var n=t.sibling;if(n!==null){n.return=t.return,I=n;break}I=t.return}}function _u(e){for(;I!==null;){var t=I;try{switch(t.tag){case 0:case 11:case 15:var n=t.return;try{Ul(4,t)}catch(c){_e(t,n,c)}break;case 1:var r=t.stateNode;if(typeof r.componentDidMount=="function"){var l=t.return;try{r.componentDidMount()}catch(c){_e(t,l,c)}}var o=t.return;try{xi(t)}catch(c){_e(t,o,c)}break;case 5:var i=t.return;try{xi(t)}catch(c){_e(t,i,c)}}}catch(c){_e(t,t.return,c)}if(t===e){I=null;break}var s=t.sibling;if(s!==null){s.return=t.return,I=s;break}I=t.return}}var Tf=Math.ceil,Bl=we.ReactCurrentDispatcher,zi=we.ReactCurrentOwner,mt=we.ReactCurrentBatchConfig,ee=0,Le=null,Ee=null,Ae=0,st=0,er=Jt(0),Pe=0,Ur=null,Cn=0,Wl=0,Ri=0,Br=null,Je=null,Li=0,tr=1/0,Ut=null,Ql=!1,Ii=null,ln=null,Vl=!1,on=null,Hl=0,Wr=0,Di=null,Kl=-1,Yl=0;function Ge(){return(ee&6)!==0?ke():Kl!==-1?Kl:Kl=ke()}function sn(e){return(e.mode&1)===0?1:(ee&2)!==0&&Ae!==0?Ae&-Ae:df.transition!==null?(Yl===0&&(Yl=ds()),Yl):(e=ie,e!==0||(e=window.event,e=e===void 0?16:ks(e.type)),e)}function Et(e,t,n,r){if(50<Wr)throw Wr=0,Di=null,Error(m(185));dr(e,n,r),((ee&2)===0||e!==Le)&&(e===Le&&((ee&2)===0&&(Wl|=n),Pe===4&&an(e,Ae)),be(e,r),n===1&&ee===0&&(t.mode&1)===0&&(tr=ke()+500,kl&&en()))}function be(e,t){var n=e.callbackNode;fc(e,t);var r=nl(e,e===Le?Ae:0);if(r===0)n!==null&&us(n),e.callbackNode=null,e.callbackPriority=0;else if(t=r&-r,e.callbackPriority!==t){if(n!=null&&us(n),t===1)e.tag===0?ff(ku.bind(null,e)):ia(ku.bind(null,e)),sf(function(){(ee&6)===0&&en()}),n=null;else{switch(ps(r)){case 1:n=po;break;case 4:n=cs;break;case 16:n=Jr;break;case 536870912:n=fs;break;default:n=Jr}n=Tu(n,wu.bind(null,e))}e.callbackPriority=t,e.callbackNode=n}}function wu(e,t){if(Kl=-1,Yl=0,(ee&6)!==0)throw Error(m(327));var n=e.callbackNode;if(nr()&&e.callbackNode!==n)return null;var r=nl(e,e===Le?Ae:0);if(r===0)return null;if((r&30)!==0||(r&e.expiredLanes)!==0||t)t=Gl(e,r);else{t=r;var l=ee;ee|=2;var o=$u();(Le!==e||Ae!==t)&&(Ut=null,tr=ke()+500,Pn(e,t));do try{Lf();break}catch(s){Su(e,s)}while(!0);Jo(),Bl.current=o,ee=l,Ee!==null?t=0:(Le=null,Ae=0,t=Pe)}if(t!==0){if(t===2&&(l=mo(e),l!==0&&(r=l,t=Oi(e,l))),t===1)throw n=Ur,Pn(e,0),an(e,r),be(e,ke()),n;if(t===6)an(e,r);else{if(l=e.current.alternate,(r&30)===0&&!zf(l)&&(t=Gl(e,r),t===2&&(o=mo(e),o!==0&&(r=o,t=Oi(e,o))),t===1))throw n=Ur,Pn(e,0),an(e,r),be(e,ke()),n;switch(e.finishedWork=l,e.finishedLanes=r,t){case 0:case 1:throw Error(m(345));case 2:Tn(e,Je,Ut);break;case 3:if(an(e,r),(r&130023424)===r&&(t=Li+500-ke(),10<t)){if(nl(e,0)!==0)break;if(l=e.suspendedLanes,(l&r)!==r){Ge(),e.pingedLanes|=e.suspendedLanes&l;break}e.timeoutHandle=Bo(Tn.bind(null,e,Je,Ut),t);break}Tn(e,Je,Ut);break;case 4:if(an(e,r),(r&4194240)===r)break;for(t=e.eventTimes,l=-1;0<r;){var i=31-_t(r);o=1<<i,i=t[i],i>l&&(l=i),r&=~o}if(r=l,r=ke()-r,r=(120>r?120:480>r?480:1080>r?1080:1920>r?1920:3e3>r?3e3:4320>r?4320:1960*Tf(r/1960))-r,10<r){e.timeoutHandle=Bo(Tn.bind(null,e,Je,Ut),r);break}Tn(e,Je,Ut);break;case 5:Tn(e,Je,Ut);break;default:throw Error(m(329))}}}return be(e,ke()),e.callbackNode===n?wu.bind(null,e):null}function Oi(e,t){var n=Br;return e.current.memoizedState.isDehydrated&&(Pn(e,t).flags|=256),e=Gl(e,t),e!==2&&(t=Je,Je=n,t!==null&&Fi(t)),e}function Fi(e){Je===null?Je=e:Je.push.apply(Je,e)}function zf(e){for(var t=e;;){if(t.flags&16384){var n=t.updateQueue;if(n!==null&&(n=n.stores,n!==null))for(var r=0;r<n.length;r++){var l=n[r],o=l.getSnapshot;l=l.value;try{if(!wt(o(),l))return!1}catch{return!1}}}if(n=t.child,t.subtreeFlags&16384&&n!==null)n.return=t,t=n;else{if(t===e)break;for(;t.sibling===null;){if(t.return===null||t.return===e)return!0;t=t.return}t.sibling.return=t.return,t=t.sibling}}return!0}function an(e,t){for(t&=~Ri,t&=~Wl,e.suspendedLanes|=t,e.pingedLanes&=~t,e=e.expirationTimes;0<t;){var n=31-_t(t),r=1<<n;e[n]=-1,t&=~r}}function ku(e){if((ee&6)!==0)throw Error(m(327));nr();var t=nl(e,0);if((t&1)===0)return be(e,ke()),null;var n=Gl(e,t);if(e.tag!==0&&n===2){var r=mo(e);r!==0&&(t=r,n=Oi(e,r))}if(n===1)throw n=Ur,Pn(e,0),an(e,t),be(e,ke()),n;if(n===6)throw Error(m(345));return e.finishedWork=e.current.alternate,e.finishedLanes=t,Tn(e,Je,Ut),be(e,ke()),null}function Mi(e,t){var n=ee;ee|=1;try{return e(t)}finally{ee=n,ee===0&&(tr=ke()+500,kl&&en())}}function xn(e){on!==null&&on.tag===0&&(ee&6)===0&&nr();var t=ee;ee|=1;var n=mt.transition,r=ie;try{if(mt.transition=null,ie=1,e)return e()}finally{ie=r,mt.transition=n,ee=t,(ee&6)===0&&en()}}function ji(){st=er.current,de(er)}function Pn(e,t){e.finishedWork=null,e.finishedLanes=0;var n=e.timeoutHandle;if(n!==-1&&(e.timeoutHandle=-1,of(n)),Ee!==null)for(n=Ee.return;n!==null;){var r=n;switch(Yo(r),r.tag){case 1:r=r.type.childContextTypes,r!=null&&_l();break;case 3:Zn(),de(Xe),de(We),ii();break;case 5:li(r);break;case 4:Zn();break;case 13:de(ge);break;case 19:de(ge);break;case 10:bo(r.type._context);break;case 22:case 23:ji()}n=n.return}if(Le=e,Ee=e=un(e.current,null),Ae=st=t,Pe=0,Ur=null,Ri=Wl=Cn=0,Je=Br=null,$n!==null){for(t=0;t<$n.length;t++)if(n=$n[t],r=n.interleaved,r!==null){n.interleaved=null;var l=r.next,o=n.pending;if(o!==null){var i=o.next;o.next=l,r.next=i}n.pending=r}$n=null}return e}function Su(e,t){do{var n=Ee;try{if(Jo(),Rl.current=Ol,Ll){for(var r=ye.memoizedState;r!==null;){var l=r.queue;l!==null&&(l.pending=null),r=r.next}Ll=!1}if(En=0,Re=xe=ye=null,Dr=!1,Or=0,zi.current=null,n===null||n.return===null){Pe=1,Ur=t,Ee=null;break}e:{var o=e,i=n.return,s=n,c=t;if(t=Ae,s.flags|=32768,c!==null&&typeof c=="object"&&typeof c.then=="function"){var g=c,S=s,$=S.tag;if((S.mode&1)===0&&($===0||$===11||$===15)){var k=S.alternate;k?(S.updateQueue=k.updateQueue,S.memoizedState=k.memoizedState,S.lanes=k.lanes):(S.updateQueue=null,S.memoizedState=null)}var z=Ya(i);if(z!==null){z.flags&=-257,Ga(z,i,s,o,t),z.mode&1&&Ka(o,g,t),t=z,c=g;var O=t.updateQueue;if(O===null){var F=new Set;F.add(c),t.updateQueue=F}else O.add(c);break e}else{if((t&1)===0){Ka(o,g,t),Ai();break e}c=Error(m(426))}}else if(ve&&s.mode&1){var Se=Ya(i);if(Se!==null){(Se.flags&65536)===0&&(Se.flags|=256),Ga(Se,i,s,o,t),qo(Jn(c,s));break e}}o=c=Jn(c,s),Pe!==4&&(Pe=2),Br===null?Br=[o]:Br.push(o),o=i;do{switch(o.tag){case 3:o.flags|=65536,t&=-t,o.lanes|=t;var v=Va(o,c,t);ga(o,v);break e;case 1:s=c;var p=o.type,h=o.stateNode;if((o.flags&128)===0&&(typeof p.getDerivedStateFromError=="function"||h!==null&&typeof h.componentDidCatch=="function"&&(ln===null||!ln.has(h)))){o.flags|=65536,t&=-t,o.lanes|=t;var N=Ha(o,s,t);ga(o,N);break e}}o=o.return}while(o!==null)}Eu(n)}catch(j){t=j,Ee===n&&n!==null&&(Ee=n=n.return);continue}break}while(!0)}function $u(){var e=Bl.current;return Bl.current=Ol,e===null?Ol:e}function Ai(){(Pe===0||Pe===3||Pe===2)&&(Pe=4),Le===null||(Cn&268435455)===0&&(Wl&268435455)===0||an(Le,Ae)}function Gl(e,t){var n=ee;ee|=2;var r=$u();(Le!==e||Ae!==t)&&(Ut=null,Pn(e,t));do try{Rf();break}catch(l){Su(e,l)}while(!0);if(Jo(),ee=n,Bl.current=r,Ee!==null)throw Error(m(261));return Le=null,Ae=0,Pe}function Rf(){for(;Ee!==null;)Nu(Ee)}function Lf(){for(;Ee!==null&&!nc();)Nu(Ee)}function Nu(e){var t=Pu(e.alternate,e,st);e.memoizedProps=e.pendingProps,t===null?Eu(e):Ee=t,zi.current=null}function Eu(e){var t=e;do{var n=t.alternate;if(e=t.return,(t.flags&32768)===0){if(n=Nf(n,t,st),n!==null){Ee=n;return}}else{if(n=Ef(n,t),n!==null){n.flags&=32767,Ee=n;return}if(e!==null)e.flags|=32768,e.subtreeFlags=0,e.deletions=null;else{Pe=6,Ee=null;return}}if(t=t.sibling,t!==null){Ee=t;return}Ee=t=e}while(t!==null);Pe===0&&(Pe=5)}function Tn(e,t,n){var r=ie,l=mt.transition;try{mt.transition=null,ie=1,If(e,t,n,r)}finally{mt.transition=l,ie=r}return null}function If(e,t,n,r){do nr();while(on!==null);if((ee&6)!==0)throw Error(m(327));n=e.finishedWork;var l=e.finishedLanes;if(n===null)return null;if(e.finishedWork=null,e.finishedLanes=0,n===e.current)throw Error(m(177));e.callbackNode=null,e.callbackPriority=0;var o=n.lanes|n.childLanes;if(dc(e,o),e===Le&&(Ee=Le=null,Ae=0),(n.subtreeFlags&2064)===0&&(n.flags&2064)===0||Vl||(Vl=!0,Tu(Jr,function(){return nr(),null})),o=(n.flags&15990)!==0,(n.subtreeFlags&15990)!==0||o){o=mt.transition,mt.transition=null;var i=ie;ie=1;var s=ee;ee|=4,zi.current=null,xf(e,n),vu(n,e),Jc(Ao),ol=!!jo,Ao=jo=null,e.current=n,Pf(n),rc(),ee=s,ie=i,mt.transition=o}else e.current=n;if(Vl&&(Vl=!1,on=e,Hl=l),o=e.pendingLanes,o===0&&(ln=null),ic(n.stateNode),be(e,ke()),t!==null)for(r=e.onRecoverableError,n=0;n<t.length;n++)l=t[n],r(l.value,{componentStack:l.stack,digest:l.digest});if(Ql)throw Ql=!1,e=Ii,Ii=null,e;return(Hl&1)!==0&&e.tag!==0&&nr(),o=e.pendingLanes,(o&1)!==0?e===Di?Wr++:(Wr=0,Di=e):Wr=0,en(),null}function nr(){if(on!==null){var e=ps(Hl),t=mt.transition,n=ie;try{if(mt.transition=null,ie=16>e?16:e,on===null)var r=!1;else{if(e=on,on=null,Hl=0,(ee&6)!==0)throw Error(m(331));var l=ee;for(ee|=4,I=e.current;I!==null;){var o=I,i=o.child;if((I.flags&16)!==0){var s=o.deletions;if(s!==null){for(var c=0;c<s.length;c++){var g=s[c];for(I=g;I!==null;){var S=I;switch(S.tag){case 0:case 11:case 15:Ar(8,S,o)}var $=S.child;if($!==null)$.return=S,I=$;else for(;I!==null;){S=I;var k=S.sibling,z=S.return;if(cu(S),S===g){I=null;break}if(k!==null){k.return=z,I=k;break}I=z}}}var O=o.alternate;if(O!==null){var F=O.child;if(F!==null){O.child=null;do{var Se=F.sibling;F.sibling=null,F=Se}while(F!==null)}}I=o}}if((o.subtreeFlags&2064)!==0&&i!==null)i.return=o,I=i;else e:for(;I!==null;){if(o=I,(o.flags&2048)!==0)switch(o.tag){case 0:case 11:case 15:Ar(9,o,o.return)}var v=o.sibling;if(v!==null){v.return=o.return,I=v;break e}I=o.return}}var p=e.current;for(I=p;I!==null;){i=I;var h=i.child;if((i.subtreeFlags&2064)!==0&&h!==null)h.return=i,I=h;else e:for(i=p;I!==null;){if(s=I,(s.flags&2048)!==0)try{switch(s.tag){case 0:case 11:case 15:Ul(9,s)}}catch(j){_e(s,s.return,j)}if(s===i){I=null;break e}var N=s.sibling;if(N!==null){N.return=s.return,I=N;break e}I=s.return}}if(ee=l,en(),xt&&typeof xt.onPostCommitFiberRoot=="function")try{xt.onPostCommitFiberRoot(br,e)}catch{}r=!0}return r}finally{ie=n,mt.transition=t}}return!1}function Cu(e,t,n){t=Jn(n,t),t=Va(e,t,1),e=nn(e,t,1),t=Ge(),e!==null&&(dr(e,1,t),be(e,t))}function _e(e,t,n){if(e.tag===3)Cu(e,e,n);else for(;t!==null;){if(t.tag===3){Cu(t,e,n);break}else if(t.tag===1){var r=t.stateNode;if(typeof t.type.getDerivedStateFromError=="function"||typeof r.componentDidCatch=="function"&&(ln===null||!ln.has(r))){e=Jn(n,e),e=Ha(t,e,1),t=nn(t,e,1),e=Ge(),t!==null&&(dr(t,1,e),be(t,e));break}}t=t.return}}function Df(e,t,n){var r=e.pingCache;r!==null&&r.delete(t),t=Ge(),e.pingedLanes|=e.suspendedLanes&n,Le===e&&(Ae&n)===n&&(Pe===4||Pe===3&&(Ae&130023424)===Ae&&500>ke()-Li?Pn(e,0):Ri|=n),be(e,t)}function xu(e,t){t===0&&((e.mode&1)===0?t=1:(t=tl,tl<<=1,(tl&130023424)===0&&(tl=4194304)));var n=Ge();e=Mt(e,t),e!==null&&(dr(e,t,n),be(e,n))}function Of(e){var t=e.memoizedState,n=0;t!==null&&(n=t.retryLane),xu(e,n)}function Ff(e,t){var n=0;switch(e.tag){case 13:var r=e.stateNode,l=e.memoizedState;l!==null&&(n=l.retryLane);break;case 19:r=e.stateNode;break;default:throw Error(m(314))}r!==null&&r.delete(t),xu(e,n)}var Pu;Pu=function(e,t,n){if(e!==null)if(e.memoizedProps!==t.pendingProps||Xe.current)Ze=!0;else{if((e.lanes&n)===0&&(t.flags&128)===0)return Ze=!1,$f(e,t,n);Ze=(e.flags&131072)!==0}else Ze=!1,ve&&(t.flags&1048576)!==0&&sa(t,$l,t.index);switch(t.lanes=0,t.tag){case 2:var r=t.type;jl(e,t),e=t.pendingProps;var l=Vn(t,We.current);qn(t,n),l=ui(null,t,r,e,l,n);var o=ci();return t.flags|=1,typeof l=="object"&&l!==null&&typeof l.render=="function"&&l.$$typeof===void 0?(t.tag=1,t.memoizedState=null,t.updateQueue=null,qe(r)?(o=!0,wl(t)):o=!1,t.memoizedState=l.state!==null&&l.state!==void 0?l.state:null,ni(t),l.updater=Fl,t.stateNode=l,l._reactInternals=t,hi(t,r,e,n),t=wi(null,t,r,!0,o,n)):(t.tag=0,ve&&o&&Ko(t),Ye(null,t,l,n),t=t.child),t;case 16:r=t.elementType;e:{switch(jl(e,t),e=t.pendingProps,l=r._init,r=l(r._payload),t.type=r,l=t.tag=jf(r),e=St(r,e),l){case 0:t=_i(null,t,r,e,n);break e;case 1:t=eu(null,t,r,e,n);break e;case 11:t=Xa(null,t,r,e,n);break e;case 14:t=qa(null,t,r,St(r.type,e),n);break e}throw Error(m(306,r,""))}return t;case 0:return r=t.type,l=t.pendingProps,l=t.elementType===r?l:St(r,l),_i(e,t,r,l,n);case 1:return r=t.type,l=t.pendingProps,l=t.elementType===r?l:St(r,l),eu(e,t,r,l,n);case 3:e:{if(tu(t),e===null)throw Error(m(387));r=t.pendingProps,o=t.memoizedState,l=o.element,ha(e,t),Tl(t,r,null,n);var i=t.memoizedState;if(r=i.element,o.isDehydrated)if(o={element:r,isDehydrated:!1,cache:i.cache,pendingSuspenseBoundaries:i.pendingSuspenseBoundaries,transitions:i.transitions},t.updateQueue.baseState=o,t.memoizedState=o,t.flags&256){l=Jn(Error(m(423)),t),t=nu(e,t,r,n,l);break e}else if(r!==l){l=Jn(Error(m(424)),t),t=nu(e,t,r,n,l);break e}else for(it=Zt(t.stateNode.containerInfo.firstChild),ot=t,ve=!0,kt=null,n=ma(t,null,r,n),t.child=n;n;)n.flags=n.flags&-3|4096,n=n.sibling;else{if(Yn(),r===l){t=At(e,t,n);break e}Ye(e,t,r,n)}t=t.child}return t;case 5:return _a(t),e===null&&Xo(t),r=t.type,l=t.pendingProps,o=e!==null?e.memoizedProps:null,i=l.children,Uo(r,l)?i=null:o!==null&&Uo(r,o)&&(t.flags|=32),ba(e,t),Ye(e,t,i,n),t.child;case 6:return e===null&&Xo(t),null;case 13:return ru(e,t,n);case 4:return ri(t,t.stateNode.containerInfo),r=t.pendingProps,e===null?t.child=Gn(t,null,r,n):Ye(e,t,r,n),t.child;case 11:return r=t.type,l=t.pendingProps,l=t.elementType===r?l:St(r,l),Xa(e,t,r,l,n);case 7:return Ye(e,t,t.pendingProps,n),t.child;case 8:return Ye(e,t,t.pendingProps.children,n),t.child;case 12:return Ye(e,t,t.pendingProps.children,n),t.child;case 10:e:{if(r=t.type._context,l=t.pendingProps,o=t.memoizedProps,i=l.value,ce(Cl,r._currentValue),r._currentValue=i,o!==null)if(wt(o.value,i)){if(o.children===l.children&&!Xe.current){t=At(e,t,n);break e}}else for(o=t.child,o!==null&&(o.return=t);o!==null;){var s=o.dependencies;if(s!==null){i=o.child;for(var c=s.firstContext;c!==null;){if(c.context===r){if(o.tag===1){c=jt(-1,n&-n),c.tag=2;var g=o.updateQueue;if(g!==null){g=g.shared;var S=g.pending;S===null?c.next=c:(c.next=S.next,S.next=c),g.pending=c}}o.lanes|=n,c=o.alternate,c!==null&&(c.lanes|=n),ei(o.return,n,t),s.lanes|=n;break}c=c.next}}else if(o.tag===10)i=o.type===t.type?null:o.child;else if(o.tag===18){if(i=o.return,i===null)throw Error(m(341));i.lanes|=n,s=i.alternate,s!==null&&(s.lanes|=n),ei(i,n,t),i=o.sibling}else i=o.child;if(i!==null)i.return=o;else for(i=o;i!==null;){if(i===t){i=null;break}if(o=i.sibling,o!==null){o.return=i.return,i=o;break}i=i.return}o=i}Ye(e,t,l.children,n),t=t.child}return t;case 9:return l=t.type,r=t.pendingProps.children,qn(t,n),l=dt(l),r=r(l),t.flags|=1,Ye(e,t,r,n),t.child;case 14:return r=t.type,l=St(r,t.pendingProps),l=St(r.type,l),qa(e,t,r,l,n);case 15:return Za(e,t,t.type,t.pendingProps,n);case 17:return r=t.type,l=t.pendingProps,l=t.elementType===r?l:St(r,l),jl(e,t),t.tag=1,qe(r)?(e=!0,wl(t)):e=!1,qn(t,n),Wa(t,r,l),hi(t,r,l,n),wi(null,t,r,!0,e,n);case 19:return ou(e,t,n);case 22:return Ja(e,t,n)}throw Error(m(156,t.tag))};function Tu(e,t){return as(e,t)}function Mf(e,t,n,r){this.tag=e,this.key=n,this.sibling=this.child=this.return=this.stateNode=this.type=this.elementType=null,this.index=0,this.ref=null,this.pendingProps=t,this.dependencies=this.memoizedState=this.updateQueue=this.memoizedProps=null,this.mode=r,this.subtreeFlags=this.flags=0,this.deletions=null,this.childLanes=this.lanes=0,this.alternate=null}function vt(e,t,n,r){return new Mf(e,t,n,r)}function Ui(e){return e=e.prototype,!(!e||!e.isReactComponent)}function jf(e){if(typeof e=="function")return Ui(e)?1:0;if(e!=null){if(e=e.$$typeof,e===Te)return 11;if(e===Be)return 14}return 2}function un(e,t){var n=e.alternate;return n===null?(n=vt(e.tag,t,e.key,e.mode),n.elementType=e.elementType,n.type=e.type,n.stateNode=e.stateNode,n.alternate=e,e.alternate=n):(n.pendingProps=t,n.type=e.type,n.flags=0,n.subtreeFlags=0,n.deletions=null),n.flags=e.flags&14680064,n.childLanes=e.childLanes,n.lanes=e.lanes,n.child=e.child,n.memoizedProps=e.memoizedProps,n.memoizedState=e.memoizedState,n.updateQueue=e.updateQueue,t=e.dependencies,n.dependencies=t===null?null:{lanes:t.lanes,firstContext:t.firstContext},n.sibling=e.sibling,n.index=e.index,n.ref=e.ref,n}function Xl(e,t,n,r,l,o){var i=2;if(r=e,typeof e=="function")Ui(e)&&(i=1);else if(typeof e=="string")i=5;else e:switch(e){case Fe:return zn(n.children,l,o,t);case Ue:i=8,l|=8;break;case Ct:return e=vt(12,n,t,l|2),e.elementType=Ct,e.lanes=o,e;case ze:return e=vt(13,n,t,l),e.elementType=ze,e.lanes=o,e;case Ce:return e=vt(19,n,t,l),e.elementType=Ce,e.lanes=o,e;case Q:return ql(n,l,o,t);default:if(typeof e=="object"&&e!==null)switch(e.$$typeof){case gt:i=10;break e;case yt:i=9;break e;case Te:i=11;break e;case Be:i=14;break e;case U:i=16,r=null;break e}throw Error(m(130,e==null?e:typeof e,""))}return t=vt(i,n,t,l),t.elementType=e,t.type=r,t.lanes=o,t}function zn(e,t,n,r){return e=vt(7,e,r,t),e.lanes=n,e}function ql(e,t,n,r){return e=vt(22,e,r,t),e.elementType=Q,e.lanes=n,e.stateNode={isHidden:!1},e}function Bi(e,t,n){return e=vt(6,e,null,t),e.lanes=n,e}function Wi(e,t,n){return t=vt(4,e.children!==null?e.children:[],e.key,t),t.lanes=n,t.stateNode={containerInfo:e.containerInfo,pendingChildren:null,implementation:e.implementation},t}function Af(e,t,n,r,l){this.tag=t,this.containerInfo=e,this.finishedWork=this.pingCache=this.current=this.pendingChildren=null,this.timeoutHandle=-1,this.callbackNode=this.pendingContext=this.context=null,this.callbackPriority=0,this.eventTimes=vo(0),this.expirationTimes=vo(-1),this.entangledLanes=this.finishedLanes=this.mutableReadLanes=this.expiredLanes=this.pingedLanes=this.suspendedLanes=this.pendingLanes=0,this.entanglements=vo(0),this.identifierPrefix=r,this.onRecoverableError=l,this.mutableSourceEagerHydrationData=null}function Qi(e,t,n,r,l,o,i,s,c){return e=new Af(e,t,n,s,c),t===1?(t=1,o===!0&&(t|=8)):t=0,o=vt(3,null,null,t),e.current=o,o.stateNode=e,o.memoizedState={element:r,isDehydrated:n,cache:null,transitions:null,pendingSuspenseBoundaries:null},ni(o),e}function Uf(e,t,n){var r=3<arguments.length&&arguments[3]!==void 0?arguments[3]:null;return{$$typeof:Oe,key:r==null?null:""+r,children:e,containerInfo:t,implementation:n}}function zu(e){if(!e)return bt;e=e._reactInternals;e:{if(yn(e)!==e||e.tag!==1)throw Error(m(170));var t=e;do{switch(t.tag){case 3:t=t.stateNode.context;break e;case 1:if(qe(t.type)){t=t.stateNode.__reactInternalMemoizedMergedChildContext;break e}}t=t.return}while(t!==null);throw Error(m(171))}if(e.tag===1){var n=e.type;if(qe(n))return la(e,n,t)}return t}function Ru(e,t,n,r,l,o,i,s,c){return e=Qi(n,r,!0,e,l,o,i,s,c),e.context=zu(null),n=e.current,r=Ge(),l=sn(n),o=jt(r,l),o.callback=t??null,nn(n,o,l),e.current.lanes=l,dr(e,l,r),be(e,r),e}function Zl(e,t,n,r){var l=t.current,o=Ge(),i=sn(l);return n=zu(n),t.context===null?t.context=n:t.pendingContext=n,t=jt(o,i),t.payload={element:e},r=r===void 0?null:r,r!==null&&(t.callback=r),e=nn(l,t,i),e!==null&&(Et(e,l,i,o),Pl(e,l,i)),i}function Jl(e){return e=e.current,e.child?(e.child.tag===5,e.child.stateNode):null}function Lu(e,t){if(e=e.memoizedState,e!==null&&e.dehydrated!==null){var n=e.retryLane;e.retryLane=n!==0&&n<t?n:t}}function Vi(e,t){Lu(e,t),(e=e.alternate)&&Lu(e,t)}function Bf(){return null}var Iu=typeof reportError=="function"?reportError:function(e){console.error(e)};function Hi(e){this._internalRoot=e}bl.prototype.render=Hi.prototype.render=function(e){var t=this._internalRoot;if(t===null)throw Error(m(409));Zl(e,t,null,null)},bl.prototype.unmount=Hi.prototype.unmount=function(){var e=this._internalRoot;if(e!==null){this._internalRoot=null;var t=e.containerInfo;xn(function(){Zl(null,e,null,null)}),t[It]=null}};function bl(e){this._internalRoot=e}bl.prototype.unstable_scheduleHydration=function(e){if(e){var t=hs();e={blockedOn:null,target:e,priority:t};for(var n=0;n<Gt.length&&t!==0&&t<Gt[n].priority;n++);Gt.splice(n,0,e),n===0&&_s(e)}};function Ki(e){return!(!e||e.nodeType!==1&&e.nodeType!==9&&e.nodeType!==11)}function eo(e){return!(!e||e.nodeType!==1&&e.nodeType!==9&&e.nodeType!==11&&(e.nodeType!==8||e.nodeValue!==" react-mount-point-unstable "))}function Du(){}function Wf(e,t,n,r,l){if(l){if(typeof r=="function"){var o=r;r=function(){var g=Jl(i);o.call(g)}}var i=Ru(t,r,e,0,null,!1,!1,"",Du);return e._reactRootContainer=i,e[It]=i.current,Cr(e.nodeType===8?e.parentNode:e),xn(),i}for(;l=e.lastChild;)e.removeChild(l);if(typeof r=="function"){var s=r;r=function(){var g=Jl(c);s.call(g)}}var c=Qi(e,0,!1,null,null,!1,!1,"",Du);return e._reactRootContainer=c,e[It]=c.current,Cr(e.nodeType===8?e.parentNode:e),xn(function(){Zl(t,c,n,r)}),c}function to(e,t,n,r,l){var o=n._reactRootContainer;if(o){var i=o;if(typeof l=="function"){var s=l;l=function(){var c=Jl(i);s.call(c)}}Zl(t,i,e,l)}else i=Wf(n,t,e,l,r);return Jl(i)}ms=function(e){switch(e.tag){case 3:var t=e.stateNode;if(t.current.memoizedState.isDehydrated){var n=fr(t.pendingLanes);n!==0&&(ho(t,n|1),be(t,ke()),(ee&6)===0&&(tr=ke()+500,en()))}break;case 13:xn(function(){var r=Mt(e,1);if(r!==null){var l=Ge();Et(r,e,1,l)}}),Vi(e,1)}},go=function(e){if(e.tag===13){var t=Mt(e,134217728);if(t!==null){var n=Ge();Et(t,e,134217728,n)}Vi(e,134217728)}},vs=function(e){if(e.tag===13){var t=sn(e),n=Mt(e,t);if(n!==null){var r=Ge();Et(n,e,t,r)}Vi(e,t)}},hs=function(){return ie},gs=function(e,t){var n=ie;try{return ie=e,t()}finally{ie=n}},ao=function(e,t,n){switch(t){case"input":if(mn(e,n),t=n.name,n.type==="radio"&&t!=null){for(n=e;n.parentNode;)n=n.parentNode;for(n=n.querySelectorAll("input[name="+JSON.stringify(""+t)+'][type="radio"]'),t=0;t<n.length;t++){var r=n[t];if(r!==e&&r.form===e.form){var l=yl(r);if(!l)throw Error(m(90));dn(r),mn(r,l)}}}break;case"textarea":Kr(e,n);break;case"select":t=n.value,t!=null&&Qt(e,!!n.multiple,t,!1)}},ts=Mi,ns=xn;var Qf={usingClientEntryPoint:!1,Events:[Tr,Wn,yl,bi,es,Mi]},Qr={findFiberByHostInstance:_n,bundleType:0,version:"18.3.1",rendererPackageName:"react-dom"},Vf={bundleType:Qr.bundleType,version:Qr.version,rendererPackageName:Qr.rendererPackageName,rendererConfig:Qr.rendererConfig,overrideHookState:null,overrideHookStateDeletePath:null,overrideHookStateRenamePath:null,overrideProps:null,overridePropsDeletePath:null,overridePropsRenamePath:null,setErrorHandler:null,setSuspenseHandler:null,scheduleUpdate:null,currentDispatcherRef:we.ReactCurrentDispatcher,findHostInstanceByFiber:function(e){return e=is(e),e===null?null:e.stateNode},findFiberByHostInstance:Qr.findFiberByHostInstance||Bf,findHostInstancesForRefresh:null,scheduleRefresh:null,scheduleRoot:null,setRefreshHandler:null,getCurrentFiber:null,reconcilerVersion:"18.3.1-next-f1338f8080-20240426"};if(typeof __REACT_DEVTOOLS_GLOBAL_HOOK__<"u"){var no=__REACT_DEVTOOLS_GLOBAL_HOOK__;if(!no.isDisabled&&no.supportsFiber)try{br=no.inject(Vf),xt=no}catch{}}return et.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED=Qf,et.createPortal=function(e,t){var n=2<arguments.length&&arguments[2]!==void 0?arguments[2]:null;if(!Ki(t))throw Error(m(200));return Uf(e,t,null,n)},et.createRoot=function(e,t){if(!Ki(e))throw Error(m(299));var n=!1,r="",l=Iu;return t!=null&&(t.unstable_strictMode===!0&&(n=!0),t.identifierPrefix!==void 0&&(r=t.identifierPrefix),t.onRecoverableError!==void 0&&(l=t.onRecoverableError)),t=Qi(e,1,!1,null,null,n,!1,r,l),e[It]=t.current,Cr(e.nodeType===8?e.parentNode:e),new Hi(t)},et.findDOMNode=function(e){if(e==null)return null;if(e.nodeType===1)return e;var t=e._reactInternals;if(t===void 0)throw typeof e.render=="function"?Error(m(188)):(e=Object.keys(e).join(","),Error(m(268,e)));return e=is(t),e=e===null?null:e.stateNode,e},et.flushSync=function(e){return xn(e)},et.hydrate=function(e,t,n){if(!eo(t))throw Error(m(200));return to(null,e,t,!0,n)},et.hydrateRoot=function(e,t,n){if(!Ki(e))throw Error(m(405));var r=n!=null&&n.hydratedSources||null,l=!1,o="",i=Iu;if(n!=null&&(n.unstable_strictMode===!0&&(l=!0),n.identifierPrefix!==void 0&&(o=n.identifierPrefix),n.onRecoverableError!==void 0&&(i=n.onRecoverableError)),t=Ru(t,null,e,1,n??null,l,!1,o,i),e[It]=t.current,Cr(e),r)for(e=0;e<r.length;e++)n=r[e],l=n._getVersion,l=l(n._source),t.mutableSourceEagerHydrationData==null?t.mutableSourceEagerHydrationData=[n,l]:t.mutableSourceEagerHydrationData.push(n,l);return new bl(t)},et.render=function(e,t,n){if(!eo(t))throw Error(m(200));return to(null,e,t,!1,n)},et.unmountComponentAtNode=function(e){if(!eo(e))throw Error(m(40));return e._reactRootContainer?(xn(function(){to(null,null,e,!1,function(){e._reactRootContainer=null,e[It]=null})}),!0):!1},et.unstable_batchedUpdates=Mi,et.unstable_renderSubtreeIntoContainer=function(e,t,n,r){if(!eo(n))throw Error(m(200));if(e==null||e._reactInternals===void 0)throw Error(m(38));return to(e,t,n,!1,r)},et.version="18.3.1-next-f1338f8080-20240426",et}var Uu;function Jf(){if(Uu)return Gi.exports;Uu=1;function f(){if(!(typeof __REACT_DEVTOOLS_GLOBAL_HOOK__>"u"||typeof __REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE!="function"))try{__REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE(f)}catch(_){console.error(_)}}return f(),Gi.exports=Zf(),Gi.exports}var Bu;function bf(){if(Bu)return ro;Bu=1;var f=Jf();return ro.createRoot=f.createRoot,ro.hydrateRoot=f.hydrateRoot,ro}var ed=bf(),Gu=function(f,_,m,T){var C;_[0]=0;for(var x=1;x<_.length;x++){var M=_[x++],G=_[x]?(_[0]|=M?1:2,m[_[x++]]):_[++x];M===3?T[0]=G:M===4?T[1]=Object.assign(T[1]||{},G):M===5?(T[1]=T[1]||{})[_[++x]]=G:M===6?T[1][_[++x]]+=G+"":M?(C=f.apply(G,Gu(f,G,m,["",null])),T.push(C),G[0]?_[0]|=2:(_[x-2]=0,_[x]=C)):T.push(G)}return T},Wu=new Map;function td(f){var _=Wu.get(this);return _||(_=new Map,Wu.set(this,_)),(_=Gu(this,_.get(f)||(_.set(f,_=(function(m){for(var T,C,x=1,M="",G="",E=[0],se=function(Z){x===1&&(Z||(M=M.replace(/^\s*\n\s*|\s*\n\s*$/g,"")))?E.push(0,Z,M):x===3&&(Z||M)?(E.push(3,Z,M),x=2):x===2&&M==="..."&&Z?E.push(4,Z,0):x===2&&M&&!Z?E.push(5,0,!0,M):x>=5&&((M||!Z&&x===5)&&(E.push(x,0,M,C),x=6),Z&&(E.push(x,Z,0,C),x=6)),M=""},he=0;he<m.length;he++){he&&(x===1&&se(),se(he));for(var ne=0;ne<m[he].length;ne++)T=m[he][ne],x===1?T==="<"?(se(),E=[E],x=3):M+=T:x===4?M==="--"&&T===">"?(x=1,M=""):M=T+M[0]:G?T===G?G="":M+=T:T==='"'||T==="'"?G=T:T===">"?(se(),x=1):x&&(T==="="?(x=5,C=M,M=""):T==="/"&&(x<5||m[he][ne+1]===">")?(se(),x===3&&(E=E[0]),x=E,(E=E[0]).push(2,0,x),x=0):T===" "||T==="	"||T===`
`||T==="\r"?(se(),x=2):M+=T),x===3&&M==="!--"&&(x=4,E=E[0])}return se(),E})(f)),_),arguments,[])).length>1?_:_[0]}const nd="https://publicfiles-envolve.s3.us-west-2.amazonaws.com/logo_black.svg",rd="https://publicfiles-envolve.s3.us-west-2.amazonaws.com/logo_white.svg",ld=7e3,od=3e3,Xu="whatsmeow-dashboard-theme",H=td.bind(Gf.createElement);function Ne(f){return encodeURIComponent(f)}function id(){const f=localStorage.getItem(Xu);return f==="light"||f==="dark"?f:window.matchMedia&&window.matchMedia("(prefers-color-scheme: dark)").matches?"dark":"light"}async function $e(f,_={}){const m={..._},T=new Headers(m.headers||{});m.body&&!(m.body instanceof FormData)&&!T.has("Content-Type")&&T.set("Content-Type","application/json");const C=await fetch(f,{...m,headers:T});if(!C.ok){let M=`${C.status} ${C.statusText}`;try{if((C.headers.get("content-type")||"").includes("application/json")){const E=await C.json();M=E?.error?.message||E?.message||E?.error||JSON.stringify(E)}else{const E=await C.text();E&&(M=E)}}catch{}throw new Error(M)}return(C.headers.get("content-type")||"").includes("application/json")?C.json():null}function sd(){return`${Date.now()}-${Math.random().toString(16).slice(2)}`}function Zi(){return{reject_calls:!1,reject_calls_message:"",mark_read_on_message:!1,ignore_status_broadcast:!1,ignore_newsletters:!1,outgoing_messages_coex:!1,always_online:!1,always_online_interval_seconds:"60",proxy_enabled:!1,proxy_protocol:"",proxy_host:"",proxy_port:"",proxy_username:"",proxy_password:""}}function Qu(f){const _=f?.config||{};return{reject_calls:!!_.reject_calls,reject_calls_message:_.reject_calls_message||"",mark_read_on_message:!!_.mark_read_on_message,ignore_status_broadcast:!!_.ignore_status_broadcast,ignore_newsletters:!!_.ignore_newsletters,outgoing_messages_coex:!!_.outgoing_messages_coex,always_online:!!_.always_online,always_online_interval_seconds:String(Number.isFinite(Number(_.always_online_interval_seconds))?Number(_.always_online_interval_seconds):60),proxy_enabled:!!_.proxy_enabled,proxy_protocol:_.proxy_protocol||"",proxy_host:_.proxy_host||"",proxy_port:_.proxy_port||"",proxy_username:_.proxy_username||"",proxy_password:_.proxy_password||""}}function rr(){return{id:"",name:"",url:"",url_absolute:"",header:"",token:"",timeout_ms:"5000",enabled:!0,append_session:!0,send_incoming_messages:!0,send_outgoing_messages:!0,send_update_messages:!0,send_group_messages:!0,send_newsletter_messages:!1}}function ad(f){return{id:f?.id||"",name:f?.name||"",url:f?.url||"",url_absolute:f?.url_absolute||"",header:f?.header||"",token:f?.token||"",timeout_ms:String(Number.isFinite(Number(f?.timeout_ms))?Number(f.timeout_ms):5e3),enabled:!!f?.enabled,append_session:!!f?.append_session,send_incoming_messages:!!f?.send_incoming_messages,send_outgoing_messages:!!f?.send_outgoing_messages,send_update_messages:!!f?.send_update_messages,send_group_messages:!!f?.send_group_messages,send_newsletter_messages:!!f?.send_newsletter_messages}}function lr(){return{id:"",name:"",language:"pt_BR",category:"",header_type:"",header_text:"",header_link:"",header_id:"",header_filename:"",body:"",footer:""}}function ud(f){const _=f?.header||{};return{id:f?.id||"",name:f?.name||"",language:f?.language||"pt_BR",category:f?.category||"",header_type:_?.type||"",header_text:_?.text||"",header_link:_?.link||"",header_id:_?.id||"",header_filename:_?.filename||"",body:f?.body||"",footer:f?.footer||""}}function cd(f){return{id:f.id.trim(),name:f.name.trim(),language:f.language.trim(),category:f.category.trim(),header:{type:f.header_type.trim(),text:f.header_text.trim(),link:f.header_link.trim(),id:f.header_id.trim(),filename:f.header_filename.trim()},body:f.body.trim(),footer:f.footer.trim()}}function fn(f){const _=String(f?.status||"").trim().toLowerCase();return _==="connected"?"online":_||(f?.connected?"online":"offline")}function qu(f){return f==="online"?"online":f==="connecting"?"connecting":"offline"}function Vu(f,_){const m=fn(f),T=fn(_),C=(f?.assigned?0:4)+(m==="online"?0:m==="connecting"?1:2),x=(_?.assigned?0:4)+(T==="online"?0:T==="connecting"?1:2);return C!==x?C-x:String(f?.session_id||"").localeCompare(String(_?.session_id||""),"pt-BR")}function fd(f){return f?new Intl.DateTimeFormat("pt-BR",{hour:"2-digit",minute:"2-digit",second:"2-digit"}).format(f):"Nunca"}function dd(f,_=56){return f?f.length>_?`${f.slice(0,_)}…`:f:"—"}function Hu(f){return String(f||"").replace(/\D+/g,"")}function pd(f){return f.key==="Enter"||f.key===" "}function lo({label:f,value:_,caption:m,tone:T="neutral"}){return H`
    <article className=${`stat-card stat-card--${T}`}>
      <p className="stat-card__label">${f}</p>
      <strong className="stat-card__value">${_}</strong>
      <span className="stat-card__caption">${m}</span>
    </article>
  `}function pe({label:f,value:_,onInput:m,placeholder:T="",type:C="text",min:x,step:M}){return H`
    <label className="field">
      <span>${f}</span>
      <input
        type=${C}
        value=${_}
        placeholder=${T}
        min=${x}
        step=${M}
        onInput=${G=>m(G.target.value)}
      />
    </label>
  `}function oo({label:f,value:_,onInput:m,placeholder:T="",rows:C=4}){return H`
    <label className="field">
      <span>${f}</span>
      <textarea
        rows=${C}
        value=${_}
        placeholder=${T}
        onInput=${x=>m(x.target.value)}
      ></textarea>
    </label>
  `}function Ke({label:f,description:_="",checked:m,onChange:T}){return H`
    <label className="check-field">
      <input
        type="checkbox"
        checked=${m}
        onChange=${C=>T(C.target.checked)}
      />
      <span>
        <strong>${f}</strong>
        ${_?H`<small>${_}</small>`:null}
      </span>
    </label>
  `}function md({active:f,label:_,count:m,onClick:T}){return H`
    <button
      type="button"
      className=${`filter-chip ${f?"filter-chip--active":""}`}
      onClick=${T}
    >
      <span>${_}</span>
      <strong>${m}</strong>
    </button>
  `}function Bt({label:f,tone:_="ghost",onClick:m,disabled:T=!1}){return H`
    <button
      type="button"
      className=${`mini-button mini-button--${_}`}
      disabled=${T}
      onClick=${C=>{C.stopPropagation(),m()}}
    >
      ${f}
    </button>
  `}function vd({session:f,selected:_,busyMap:m,onSelect:T,onSessionAction:C,onOpenTab:x,onShowQR:M}){const G=fn(f),E=qu(G),se=f.assigned?"Desatribuir":"Atribuir",he=f.assigned?"unassign":"assign";return H`
    <article
      className=${`session-card ${_?"session-card--selected":""}`}
      role="button"
      tabIndex="0"
      onClick=${()=>T(f.session_id)}
      onKeyDown=${ne=>{pd(ne)&&(ne.preventDefault(),T(f.session_id))}}
    >
      <div className="session-card__top">
        <div>
          <p className="session-card__eyebrow">Instância</p>
          <h3 className="session-card__title">${f.session_id}</h3>
        </div>
        <span className=${`status-pill status-pill--${E}`}>${G}</span>
      </div>

      <div className="tag-row">
        <span className="tag">${f.connected?"Conectada":"Desconectada"}</span>
        <span className="tag">${f.assigned?"Atribuída":"Sem atribuição"}</span>
        <span className="tag">${f.has_qr?"QR disponível":"Sem QR"}</span>
      </div>

      <div className="session-card__meta">
        <span>Login: ${f.logged_in?"ativo":"pendente"}</span>
        <span>Device: ${f.device_jid?dd(f.device_jid,30):"—"}</span>
        <span>WA Web: ${f.whatsapp_web_version||"—"}</span>
      </div>

      <div className="session-card__actions">
        <${Bt}
          label="QR"
          onClick=${()=>M(f.session_id)}
          disabled=${!f.has_qr}
        />
        <${Bt}
          label="Conectar"
          onClick=${()=>C("connect",f.session_id)}
          disabled=${!!m[`connect:${f.session_id}`]}
        />
        <${Bt}
          label="Recarregar"
          onClick=${()=>C("reload",f.session_id)}
          disabled=${!!m[`reload:${f.session_id}`]}
        />
        <${Bt}
          label="Desconectar"
          tone="danger"
          onClick=${()=>C("disconnect",f.session_id)}
          disabled=${!!m[`disconnect:${f.session_id}`]}
        />
        <${Bt}
          label="Excluir"
          tone="danger"
          onClick=${()=>C("delete",f.session_id)}
          disabled=${!!m[`delete:${f.session_id}`]}
        />
        <${Bt}
          label=${se}
          onClick=${()=>C(he,f.session_id)}
          disabled=${!!m[`${he}:${f.session_id}`]}
        />
        <${Bt}
          label="Config"
          onClick=${()=>x(f.session_id,"config")}
        />
        <${Bt}
          label="Hooks"
          onClick=${()=>x(f.session_id,"webhooks")}
        />
        <${Bt}
          label="Templates"
          onClick=${()=>x(f.session_id,"templates")}
        />
      </div>
    </article>
  `}function at({label:f,value:_,code:m=!1}){return H`
    <div className="detail-item">
      <span>${f}</span>
      ${m?H`<code>${_}</code>`:H`<strong>${_}</strong>`}
    </div>
  `}function hd({sessionId:f,onClose:_}){const[m,T]=le.useState(Date.now());return le.useEffect(()=>{if(!f)return;T(Date.now());const C=window.setInterval(()=>T(Date.now()),od),x=M=>{M.key==="Escape"&&_()};return window.addEventListener("keydown",x),()=>{window.clearInterval(C),window.removeEventListener("keydown",x)}},[f]),f?H`
    <div
      className="modal-backdrop"
      onClick=${C=>{C.target===C.currentTarget&&_()}}
    >
      <div className="modal-card">
        <div className="modal-card__header">
          <div>
            <p className="eyebrow">Pareamento</p>
            <h3>QR da sessão ${f}</h3>
          </div>
          <button type="button" className="ghost-button" onClick=${_}>Fechar</button>
        </div>
        <img
          className="modal-card__image"
          src=${`/sessions/${Ne(f)}/qr?ts=${m}`}
          alt=${`QR da sessão ${f}`}
        />
        <p className="hint-text">
          O QR atualiza automaticamente a cada poucos segundos enquanto a sessão estiver em pareamento.
        </p>
      </div>
    </div>
  `:null}function gd({item:f,onEdit:_,onDelete:m}){return H`
    <article className="webhook-card">
      <div className="webhook-card__top">
        <div>
          <strong>${f.id}</strong>
          <p>${f.name||"Webhook sem nome"}</p>
        </div>
        <div className="inline-actions">
          <button type="button" className="text-button" onClick=${_}>Editar</button>
          <button type="button" className="text-button text-button--danger" onClick=${m}>
            Excluir
          </button>
        </div>
      </div>

      <code className="webhook-card__url">${f.url_absolute||f.url||"—"}</code>

      <div className="tag-row">
        <span className="tag">${f.enabled?"Ativo":"Desligado"}</span>
        <span className="tag">${f.append_session?"Anexa sessão":"URL fixa"}</span>
        <span className="tag">${f.send_incoming_messages?"Incoming":"Sem incoming"}</span>
        <span className="tag">${f.send_outgoing_messages?"Outgoing":"Sem outgoing"}</span>
        <span className="tag">${f.send_update_messages?"Update":"Sem update"}</span>
      </div>
    </article>
  `}function yd({item:f,onEdit:_,onDelete:m}){const T=f?.header||{},C=T.type||(T.text?"text":"sem cabeçalho");return H`
    <article className="webhook-card">
      <div className="webhook-card__top">
        <div>
          <strong>${f.name}</strong>
          <p>${f.language||"idioma padrão"}${f.category?` · ${f.category}`:""}</p>
        </div>
        <div className="inline-actions">
          <button type="button" className="text-button" onClick=${_}>Editar</button>
          <button type="button" className="text-button text-button--danger" onClick=${m}>
            Excluir
          </button>
        </div>
      </div>

      <code className="webhook-card__url">${f.body||"—"}</code>

      <div className="tag-row">
        <span className="tag">Header: ${C}</span>
        <span className="tag">${f.footer?"Com footer":"Sem footer"}</span>
        <span className="tag">${f.id}</span>
      </div>
    </article>
  `}function _d(){const[f,_]=le.useState(id),[m,T]=le.useState(null),[C,x]=le.useState([]),[M,G]=le.useState("home"),[E,se]=le.useState(""),[he,ne]=le.useState("quick"),[Z,ut]=le.useState(!0),[De,ue]=le.useState(!0),[re,ht]=le.useState(null),[tt,nt]=le.useState(""),[we,rt]=le.useState("all"),[Oe,Fe]=le.useState([]),[Ue,Ct]=le.useState({}),[gt,yt]=le.useState(""),[Te,ze]=le.useState({id:"",assign:!0,connect:!0}),[Ce,Be]=le.useState({session:"",to:"",body:""}),[U,Q]=le.useState({loading:!1,saving:!1,snapshot:null,form:Zi()}),[P,A]=le.useState({phone:"",code:"",loading:!1}),[w,d]=le.useState({loading:!1,saving:!1,items:[],form:rr()}),[y,D]=le.useState({loading:!1,saving:!1,items:[],form:lr()}),Y=le.useRef(null);function R(a,u="info"){const L=sd();Fe(K=>[...K,{id:L,message:a,tone:u}]),window.setTimeout(()=>{Fe(K=>K.filter(me=>me.id!==L))},3400)}function X(a,u){Ct(L=>{const K={...L};return u?K[a]=!0:delete K[a],K})}async function J({silent:a=!1}={}){a||ut(!0);try{const[u,L]=await Promise.all([$e("/api/v1/worker"),$e("/api/v1/sessions")]);T(u),x(Array.isArray(L?.sessions)?L.sessions:[]),ht(new Date)}catch(u){R(u.message,"error")}finally{a||ut(!1)}}Y.current=J;async function b(a){const u=String(a||"").trim();if(u){Q(L=>({...L,loading:!0}));try{const L=await $e(`/api/v1/sessions/${Ne(u)}/config`);Q(K=>({...K,loading:!1,snapshot:L,form:Qu(L)}))}catch(L){Q(K=>({...K,loading:!1})),R(L.message,"error")}}}async function ae(a,{resetForm:u=!1}={}){const L=String(a||"").trim();if(L){d(K=>({...K,loading:!0}));try{const K=await $e(`/api/v1/sessions/${Ne(L)}/webhooks`);d(me=>({...me,loading:!1,items:Array.isArray(K?.webhooks)?K.webhooks:[],form:u?rr():me.form}))}catch(K){d(me=>({...me,loading:!1})),R(K.message,"error")}}}async function Me(a,{resetForm:u=!1}={}){const L=String(a||"").trim();if(L){D(K=>({...K,loading:!0}));try{const K=await $e(`/api/v1/sessions/${Ne(L)}/templates`);D(me=>({...me,loading:!1,items:Array.isArray(K?.templates)?K.templates:[],form:u?lr():me.form}))}catch(K){D(me=>({...me,loading:!1})),R(K.message,"error")}}}le.useEffect(()=>{document.documentElement.dataset.theme=f,localStorage.setItem(Xu,f)},[f]),le.useEffect(()=>{document.title=m?.ui_title||"WhatsMeow Control Center"},[m]),le.useEffect(()=>{J()},[]),le.useEffect(()=>{if(!De)return;const a=window.setInterval(()=>{!document.hidden&&Y.current&&Y.current({silent:!0})},ld);return()=>window.clearInterval(a)},[De]),le.useEffect(()=>{if(C.length===0){se("");return}if(E&&C.some(u=>u.session_id===E))return;const a=[...C].sort(Vu)[0];a&&se(a.session_id)},[C,E]),le.useEffect(()=>{if(!E){Q({loading:!1,saving:!1,snapshot:null,form:Zi()}),A({phone:"",code:"",loading:!1}),d({loading:!1,saving:!1,items:[],form:rr()}),D({loading:!1,saving:!1,items:[],form:lr()});return}Be(u=>({...u,session:E})),A({phone:Hu(E),code:"",loading:!1});let a=!0;return Q({loading:!0,saving:!1,snapshot:null,form:Zi()}),d({loading:!0,saving:!1,items:[],form:rr()}),D({loading:!0,saving:!1,items:[],form:lr()}),Promise.all([$e(`/api/v1/sessions/${Ne(E)}/config`),$e(`/api/v1/sessions/${Ne(E)}/webhooks`),$e(`/api/v1/sessions/${Ne(E)}/templates`)]).then(([u,L,K])=>{a&&(Q({loading:!1,saving:!1,snapshot:u,form:Qu(u)}),d({loading:!1,saving:!1,items:Array.isArray(L?.webhooks)?L.webhooks:[],form:rr()}),D({loading:!1,saving:!1,items:Array.isArray(K?.templates)?K.templates:[],form:lr()}))}).catch(u=>{a&&(Q(L=>({...L,loading:!1})),d(L=>({...L,loading:!1})),D(L=>({...L,loading:!1})),R(u.message,"error"))}),()=>{a=!1}},[E]),m?.ui_title;const oe=C.find(a=>a.session_id===E)||null,dn=C.filter(a=>a.connected||fn(a)==="online").length,pn=C.filter(a=>a.has_qr).length,Wt=C.filter(a=>a.assigned).length,Vr=[...C].sort(Vu),or=tt.trim().toLowerCase(),mn=Vr.filter(a=>we==="assigned"&&!a.assigned||we==="online"&&!(a.connected||fn(a)==="online")||we==="qr"&&!a.has_qr?!1:or?[a.session_id,a.device_jid,fn(a),a.assigned?"atribuida":"livre"].filter(Boolean).join(" ").toLowerCase().includes(or):!0),Rn=oe?fn(oe):"",ir=!!(oe&&oe.has_qr&&Rn==="connecting");async function vn(a){a.preventDefault();const u=Te.id.trim();if(!u){R("Informe o ID da sessão.","error");return}const L=`create:${u}`;X(L,!0);try{await $e("/api/v1/sessions",{method:"POST",body:JSON.stringify({id:u,assign:Te.assign,connect:Te.connect})}),R(`Sessão ${u} solicitada com sucesso.`,"success"),ze({id:"",assign:!0,connect:!0}),se(u),ne("overview"),await J({silent:!0})}catch(K){R(K.message,"error")}finally{X(L,!1)}}async function Qt(a){a.preventDefault();const u=Ce.session.trim(),L=Ce.to.trim(),K=Ce.body.trim();if(!u||!L||!K){R("Preencha sessão, destino e mensagem.","error");return}const me=`send:${u}`;X(me,!0);try{await $e(`/api/v1/sessions/${Ne(u)}/messages`,{method:"POST",body:JSON.stringify({type:"text",to:L,text:{body:K}})}),R(`Mensagem enviada pela sessão ${u}.`,"success"),Be(Vt=>({...Vt,body:""}))}catch(Vt){R(Vt.message,"error")}finally{X(me,!1)}}async function Lt(a,u){const L=String(u||"").trim();if(!L)return;if(a==="disconnect"||a==="delete"){const me=a==="delete"?`Excluir permanentemente a sessão "${L}"? Isso remove o diretório da sessão, a atribuição, overrides e webhooks.`:`Desconectar a sessão "${L}"?`;if(!window.confirm(me))return}const K=`${a}:${L}`;X(K,!0);try{a==="assign"?await $e("/api/v1/worker/assignments",{method:"POST",body:JSON.stringify({session_id:L})}):a==="unassign"?await $e(`/api/v1/worker/assignments/${Ne(L)}`,{method:"DELETE"}):a==="delete"?await $e(`/api/v1/sessions/${Ne(L)}`,{method:"DELETE"}):a==="qr"?yt(L):await $e(`/sessions/${Ne(L)}/${a}`,{method:"POST"}),a!=="qr"&&(R(a==="delete"?`Sessão ${L} excluída.`:`Ação "${a}" aplicada em ${L}.`,"success"),await J({silent:!0}))}catch(me){R(me.message,"error")}finally{X(K,!1)}}function Hr(a,u){se(a),G("instances"),ne(u)}async function Kr(a){if(a.preventDefault(),!E)return;const u=parseInt(U.form.always_online_interval_seconds,10);if(!Number.isFinite(u)||u<=0){R("O intervalo de always-online deve ser maior que zero.","error");return}if(U.form.proxy_enabled){if(!U.form.proxy_host.trim()){R("Informe o host do proxy.","error");return}if(!U.form.proxy_port.trim()){R("Informe a porta do proxy.","error");return}}Q(L=>({...L,saving:!0}));try{await $e(`/api/v1/sessions/${Ne(E)}/config`,{method:"PUT",body:JSON.stringify({reject_calls:U.form.reject_calls,reject_calls_message:U.form.reject_calls_message.trim(),mark_read_on_message:U.form.mark_read_on_message,ignore_status_broadcast:U.form.ignore_status_broadcast,ignore_newsletters:U.form.ignore_newsletters,outgoing_messages_coex:U.form.outgoing_messages_coex,always_online:U.form.always_online,always_online_interval_seconds:u,proxy_enabled:U.form.proxy_enabled,proxy_protocol:U.form.proxy_protocol.trim(),proxy_host:U.form.proxy_host.trim(),proxy_port:U.form.proxy_port.trim(),proxy_username:U.form.proxy_username.trim(),proxy_password:U.form.proxy_password})}),R(`Configuração da sessão ${E} salva.`,"success"),await b(E)}catch(L){R(L.message,"error")}finally{Q(L=>({...L,saving:!1}))}}async function Yr(a){if(a.preventDefault(),!E)return;const u=Hu(P.phone);if(!u){R("Informe o telefone para gerar o código de pareamento.","error");return}A(L=>({...L,phone:u,code:"",loading:!0}));try{const L=await $e(`/api/v1/sessions/${Ne(E)}/pairing-code`,{method:"POST",body:JSON.stringify({phone:u})}),K=String(L?.pairing_code||"").trim();if(!K)throw new Error("pairing code not returned by server");A(me=>({...me,phone:u,code:K,loading:!1})),R(`Código de pareamento gerado para ${E}.`,"success")}catch(L){A(K=>({...K,loading:!1})),R(L.message,"error")}}async function Gr(){if(E&&window.confirm(`Remover o override da sessão "${E}"?`)){Q(a=>({...a,saving:!0}));try{await $e(`/api/v1/sessions/${Ne(E)}/config`,{method:"DELETE"}),R(`Override da sessão ${E} removido.`,"success"),await b(E)}catch(a){R(a.message,"error")}finally{Q(a=>({...a,saving:!1}))}}}async function sr(a){if(a.preventDefault(),!E)return;const u=parseInt(w.form.timeout_ms,10);if(!w.form.url.trim()&&!w.form.url_absolute.trim()){R("Informe a URL base ou a URL absoluta do webhook.","error");return}if(!Number.isFinite(u)||u<=0){R("O timeout do webhook deve ser maior que zero.","error");return}d(L=>({...L,saving:!0}));try{await $e(`/api/v1/sessions/${Ne(E)}/webhooks`,{method:"POST",body:JSON.stringify({id:w.form.id.trim(),name:w.form.name.trim(),url:w.form.url.trim(),url_absolute:w.form.url_absolute.trim(),header:w.form.header.trim(),token:w.form.token.trim(),timeout_ms:u,enabled:w.form.enabled,append_session:w.form.append_session,send_incoming_messages:w.form.send_incoming_messages,send_outgoing_messages:w.form.send_outgoing_messages,send_update_messages:w.form.send_update_messages,send_group_messages:w.form.send_group_messages,send_newsletter_messages:w.form.send_newsletter_messages})}),R(`Webhook salvo para a sessão ${E}.`,"success"),await ae(E,{resetForm:!0})}catch(L){R(L.message,"error")}finally{d(L=>({...L,saving:!1}))}}async function Ln(a){if(!(!E||!a)&&window.confirm(`Excluir o webhook "${a}" da sessão "${E}"?`)){d(u=>({...u,saving:!0}));try{await $e(`/api/v1/sessions/${Ne(E)}/webhooks/${Ne(a)}`,{method:"DELETE"}),R(`Webhook ${a} excluído.`,"success"),await ae(E,{resetForm:!0})}catch(u){R(u.message,"error")}finally{d(u=>({...u,saving:!1}))}}}async function Xr(a){if(a.preventDefault(),!!E){if(!y.form.name.trim()){R("Informe o nome do template.","error");return}if(!y.form.body.trim()){R("Informe o corpo do template.","error");return}D(u=>({...u,saving:!0}));try{await $e(`/api/v1/sessions/${Ne(E)}/templates`,{method:"POST",body:JSON.stringify(cd(y.form))}),R(`Template salvo para a sessão ${E}.`,"success"),await Me(E,{resetForm:!0})}catch(u){R(u.message,"error")}finally{D(u=>({...u,saving:!1}))}}}async function hn(a){if(!(!E||!a)&&window.confirm(`Excluir o template "${a}" da sessão "${E}"?`)){D(u=>({...u,saving:!0}));try{await $e(`/api/v1/sessions/${Ne(E)}/templates/${Ne(a)}`,{method:"DELETE"}),R(`Template ${a} excluído.`,"success"),await Me(E,{resetForm:!0})}catch(u){R(u.message,"error")}finally{D(u=>({...u,saving:!1}))}}}const gn=[{id:"all",label:"Todas",count:C.length},{id:"assigned",label:"Atribuídas",count:Wt},{id:"online",label:"Online",count:dn},{id:"qr",label:"Com QR",count:pn}];return H`
    <div className="page">
      <main className="dashboard-shell">
        <aside className="system-sidebar">
          <div className="system-sidebar__brand">
            <img
              className="system-sidebar__logo"
              src=${f==="dark"?rd:nd}
              alt="Envolve"
            />
            <div>
              <strong>WhatsMeow</strong>
              <span>${m?.worker_id||"worker"}</span>
            </div>
          </div>

          <nav className="system-nav" aria-label="Navegação principal">
            <button
              type="button"
              className=${`system-nav__item ${M==="home"?"system-nav__item--active":""}`}
              onClick=${()=>G("home")}
            >
              Home
            </button>
            <button
              type="button"
              className=${`system-nav__item ${M==="instances"?"system-nav__item--active":""}`}
              onClick=${()=>G("instances")}
            >
              Instâncias
            </button>
          </nav>

          <div className="system-sidebar__meta">
            <span>${dn} online</span>
            <span>${Wt} atribuída(s)</span>
          </div>
        </aside>

        <section className="provider-topbar panel">
          <div className="provider-topbar__brand">
            <div>
              <p className="eyebrow">Provider WhatsMeow</p>
              <h1>${M==="home"?"Home":"Instâncias"}</h1>
            </div>
          </div>

          <div className="provider-topbar__actions">
            <span className="hero-panel__timestamp">
              ${re?`Atualizado às ${fd(re)}`:"Carregando dados"}
            </span>
            <button
              type="button"
              className=${`ghost-button ${De?"ghost-button--active":""}`}
              onClick=${()=>ue(a=>!a)}
            >
              ${De?"Auto ligado":"Auto desligado"}
            </button>
            <button
              type="button"
              className="solid-button"
              disabled=${Z}
              onClick=${()=>J()}
            >
              ${Z?"Atualizando...":"Atualizar"}
            </button>
            <button
              type="button"
              className="ghost-button"
              onClick=${()=>_(a=>a==="dark"?"light":"dark")}
            >
              ${f==="dark"?"Tema claro":"Tema escuro"}
            </button>
          </div>
        </section>

        ${M==="home"?H`
        <section className="summary-grid">
          <${lo}
            label="Worker ativo"
            value=${m?.worker_id||"—"}
            caption=${m?.worker_strict?"Modo strict habilitado":"Aceita sessões fora da lista"}
            tone="brand"
          />
          <${lo}
            label="Sessões atribuídas"
            value=${String(Wt)}
            caption="visíveis para este worker"
            tone="accent"
          />
          <${lo}
            label="Sessões online"
            value=${String(dn)}
            caption="conexões ativas agora"
            tone="online"
          />
          <${lo}
            label="Fila AMQP"
            value=${m?.amqp_queue||"—"}
            caption=${m?.amqp_enabled?"consumidor habilitado":"AMQP desabilitado"}
            tone="secondary"
          />
        </section>

        <section className="home-grid">
          <section className="panel">
            <div className="panel__header">
              <div>
                <p className="eyebrow">Instâncias</p>
                <h2>Resumo operacional</h2>
              </div>
            </div>
            <div className="detail-grid">
              <${at} label="Total" value=${String(C.length)} />
              <${at} label="Online" value=${String(dn)} />
              <${at} label="Com QR" value=${String(pn)} />
              <${at} label="Atribuídas" value=${String(Wt)} />
            </div>
            <div className="inline-actions inline-actions--wrap">
              <button type="button" className="solid-button" onClick=${()=>G("instances")}>
                Abrir instâncias
              </button>
              <button type="button" className="ghost-button" onClick=${()=>J()}>
                Atualizar dados
              </button>
            </div>
          </section>

          <section className="panel">
            <div className="panel__header">
              <div>
                <p className="eyebrow">Worker</p>
                <h2>Filas e bindings</h2>
              </div>
            </div>
            <div className="meta-stack">
              <div className="detail-item">
                <span>Exchange</span>
                <strong>${m?.amqp_exchange||"—"}</strong>
              </div>
              <div className="detail-item">
                <span>Queue</span>
                <strong>${m?.amqp_queue||"—"}</strong>
              </div>
              <div className="detail-item">
                <span>Auto-connect</span>
                <strong>${m?.worker_auto_connect_assigned?"sim":"não"}</strong>
              </div>
            </div>
            <div className="tag-row tag-row--spaced">
              ${(m?.amqp_bindings||[]).length?m.amqp_bindings.map(a=>H`<span key=${a} className="tag">${a}</span>`):H`<span className="empty-inline">Nenhum binding configurado.</span>`}
            </div>
          </section>
        </section>
          `:null}

        ${M==="instances"?H`
        <section className="workspace-grid">
          <div className="workspace-column workspace-column--sessions">
            <section className="panel">
              <div className="panel__header panel__header--space-between">
                <div>
                  <p className="eyebrow">Instâncias</p>
                  <h2>Instâncias</h2>
                </div>
                <span className="panel__hint">
                  ${mn.length} de ${C.length}
                </span>
              </div>

              <details className="create-instance">
                <summary>Nova instância</summary>
                <form className="stack-form" onSubmit=${vn}>
                  <${pe}
                    label="ID da sessão"
                    value=${Te.id}
                    placeholder="Ex: 5562999999999"
                    onInput=${a=>ze(u=>({...u,id:a}))}
                  />
                  <div className="check-grid check-grid--single">
                    <${Ke}
                      label="Atribuir ao worker atual"
                      checked=${Te.assign}
                      onChange=${a=>ze(u=>({...u,assign:a}))}
                    />
                    <${Ke}
                      label="Conectar imediatamente"
                      checked=${Te.connect}
                      onChange=${a=>ze(u=>({...u,connect:a}))}
                    />
                  </div>
                  <button
                    type="submit"
                    className="solid-button solid-button--wide"
                    disabled=${!!Ue[`create:${Te.id.trim()}`]}
                  >
                    Criar instância
                  </button>
                </form>
              </details>

              <div className="toolbar">
                <label className="search-field">
                  <span>Buscar</span>
                  <input
                    type="text"
                    value=${tt}
                    placeholder="Filtrar por sessão, status ou device"
                    onInput=${a=>nt(a.target.value)}
                  />
                </label>
                <div className="filter-row">
                  ${gn.map(a=>H`
                      <${md}
                        key=${a.id}
                        active=${we===a.id}
                        label=${a.label}
                        count=${a.count}
                        onClick=${()=>rt(a.id)}
                      />
                    `)}
                </div>
              </div>

              ${Z?H`
                    <div className="empty-state">
                      <strong>Carregando sessões…</strong>
                      <p>Consultando o worker e o catálogo de sessões disponíveis.</p>
                    </div>
                  `:mn.length===0?H`
                      <div className="empty-state">
                        <strong>Nenhuma sessão encontrada</strong>
                        <p>Ajuste o filtro ou crie uma nova instância.</p>
                      </div>
                    `:H`
                      <div className="session-grid">
                        ${mn.map(a=>H`
                            <${vd}
                              key=${a.session_id}
                              session=${a}
                              selected=${E===a.session_id}
                              busyMap=${Ue}
                              onSelect=${u=>{se(u),G("instances")}}
                              onSessionAction=${Lt}
                              onOpenTab=${Hr}
                              onShowQR=${yt}
                            />
                          `)}
                      </div>
                    `}
            </section>
          </div>

          <div className="workspace-column workspace-column--inspector">
            <section className="panel panel--sticky">
              ${oe?H`
                    <div className="panel__header">
                      <div>
                        <p className="eyebrow">Instância</p>
                        <h2>${oe.session_id}</h2>
                      </div>
                      <div className="session-header-actions">
                        <span className=${`status-pill status-pill--${qu(Rn)}`}>
                          ${Rn}
                        </span>
                        <button
                          type="button"
                          className="ghost-button"
                          onClick=${()=>Lt("connect",oe.session_id)}
                        >
                          Conectar
                        </button>
                        <button
                          type="button"
                          className="ghost-button"
                          onClick=${()=>Lt("reload",oe.session_id)}
                        >
                          Recarregar
                        </button>
                        <button
                          type="button"
                          className="ghost-button ghost-button--danger"
                          onClick=${()=>Lt("disconnect",oe.session_id)}
                        >
                          Desconectar
                        </button>
                        <button
                          type="button"
                          className="ghost-button ghost-button--danger"
                          onClick=${()=>Lt("delete",oe.session_id)}
                        >
                          Excluir
                        </button>
                      </div>
                    </div>

                    <div className="session-status-strip">
                      <${at} label="Worker" value=${oe.assigned?m?.worker_id||"atribuída":"livre"} />
                      <${at} label="Conexão" value=${oe.connected?"conectada":"desconectada"} />
                      <${at} label="Login" value=${oe.logged_in?"ativo":"pendente"} />
                      <${at} label="WA Web" value=${oe.whatsapp_web_version||"—"} code=${!!oe.whatsapp_web_version} />
                    </div>

                    ${ir?H`
                          <div className="inline-qr">
                            <img
                              src=${`/sessions/${Ne(oe.session_id)}/qr?ts=${Date.now()}`}
                              alt=${`QR da sessão ${oe.session_id}`}
                            />
                            <div>
                              <strong>QR disponível</strong>
                              <span>Escaneie para concluir a conexão desta instância.</span>
                            </div>
                          </div>
                        `:null}

                    <div className="tab-row">
                      ${["quick","webhooks","templates"].map(a=>H`
                          <button
                            key=${a}
                            type="button"
                            className=${`tab-button ${he===a?"tab-button--active":""}`}
                            onClick=${()=>ne(a)}
                          >
                            ${a==="quick"?"Teste rápido":a==="webhooks"?"Webhook":"Templates"}
                          </button>
                        `)}
                    </div>

                    ${he==="quick"?H`
                          <div className="detail-grid">
                            <${at}
                              label="Device JID"
                              value=${oe.device_jid||"—"}
                              code=${!!oe.device_jid}
                            />
                            <${at}
                              label="WA Web"
                              value=${oe.whatsapp_web_version||"—"}
                              code=${!!oe.whatsapp_web_version}
                            />
                            <${at}
                              label="Status"
                              value=${fn(oe)}
                            />
                            <${at}
                              label="Worker"
                              value=${oe.assigned?m?.worker_id||"atribuída":"livre"}
                            />
                            <${at}
                              label="Pairing"
                              value=${oe.has_qr?"QR disponível":"sem QR"}
                            />
                          </div>

                          <form className="stack-form" onSubmit=${Qt}>
                            <${pe}
                              label="Sessão"
                              value=${Ce.session}
                              placeholder="Selecione uma sessão ou informe manualmente"
                              onInput=${a=>Be(u=>({...u,session:a}))}
                            />
                            <${pe}
                              label="Destino"
                              value=${Ce.to}
                              placeholder="5562999999999"
                              onInput=${a=>Be(u=>({...u,to:a}))}
                            />
                            <${oo}
                              label="Mensagem"
                              value=${Ce.body}
                              rows="4"
                              placeholder="Texto da mensagem"
                              onInput=${a=>Be(u=>({...u,body:a}))}
                            />
                            <button
                              type="submit"
                              className="solid-button"
                              disabled=${!!Ue[`send:${Ce.session.trim()}`]}
                            >
                              Enviar texto
                            </button>
                          </form>

                          <form className="stack-form" onSubmit=${Yr}>
                            <div className="callout callout--muted">
                              <strong>Pareamento por código</strong>
                              <p>
                                Além do QR code, você pode gerar um código de pareamento para vincular
                                a sessão sem leitura de QR. O provider conecta a sessão, espera a fase
                                de QR internamente e então pede o código ao WhatsApp.
                              </p>
                            </div>

                            <${pe}
                              label="Telefone para o pareamento"
                              value=${P.phone}
                              placeholder="5562999999999"
                              onInput=${a=>A(u=>({...u,phone:a,code:""}))}
                            />

                            <div className="inline-actions inline-actions--wrap">
                              <button
                                type="submit"
                                className="solid-button"
                                disabled=${P.loading}
                              >
                                ${P.loading?"Gerando código...":"Gerar código de pareamento"}
                              </button>
                            </div>

                            ${P.code?H`
                                  <div className="pairing-code-card">
                                    <span className="pairing-code-card__label">Código atual</span>
                                    <code className="pairing-code-card__value">${P.code}</code>
                                    <p className="hint-text pairing-code-card__hint">
                                      Esse código expira rápido. Se o pareamento não concluir a tempo,
                                      gere um novo código.
                                    </p>
                                  </div>
                                `:null}
                          </form>
                        `:null}

                    ${he==="config"?H`
                          <form className="stack-form" onSubmit=${Kr}>
                            <div
                              className=${`callout ${U.snapshot?.has_override?"callout--accent":"callout--muted"}`}
                            >
                              <strong>
                                ${U.snapshot?.has_override?"Existe um override ativo para esta sessão.":"A sessão está usando os defaults do worker."}
                              </strong>
                              <p>
                                ${U.loading?"Carregando configuração efetiva…":"Edite os parâmetros operacionais e salve no Redis/arquivo da sessão."}
                              </p>
                            </div>

                            <div className="check-grid">
                              <${Ke}
                                label="Rejeitar ligações"
                                checked=${U.form.reject_calls}
                                onChange=${a=>Q(u=>({...u,form:{...u.form,reject_calls:a}}))}
                              />
                              <${Ke}
                                label="Marcar conversa como lida"
                                checked=${U.form.mark_read_on_message}
                                onChange=${a=>Q(u=>({...u,form:{...u.form,mark_read_on_message:a}}))}
                              />
                              <${Ke}
                                label="Ignorar status broadcast"
                                checked=${U.form.ignore_status_broadcast}
                                onChange=${a=>Q(u=>({...u,form:{...u.form,ignore_status_broadcast:a}}))}
                              />
                              <${Ke}
                                label="Ignorar newsletters"
                                checked=${U.form.ignore_newsletters}
                                onChange=${a=>Q(u=>({...u,form:{...u.form,ignore_newsletters:a}}))}
                              />
                              <${Ke}
                                label="Outgoing em smb_message_echoes"
                                checked=${U.form.outgoing_messages_coex}
                                onChange=${a=>Q(u=>({...u,form:{...u.form,outgoing_messages_coex:a}}))}
                              />
                              <${Ke}
                                label="Manter sempre online"
                                checked=${U.form.always_online}
                                onChange=${a=>Q(u=>({...u,form:{...u.form,always_online:a}}))}
                              />
                              <${Ke}
                                label="Usar proxy nesta sessão"
                                checked=${U.form.proxy_enabled}
                                onChange=${a=>Q(u=>({...u,form:{...u.form,proxy_enabled:a}}))}
                              />
                            </div>

                            <${oo}
                              label="Mensagem de rejeição"
                              value=${U.form.reject_calls_message}
                              rows="3"
                              placeholder="Opcional"
                              onInput=${a=>Q(u=>({...u,form:{...u.form,reject_calls_message:a}}))}
                            />
                            <${pe}
                              label="Intervalo always-online (segundos)"
                              value=${U.form.always_online_interval_seconds}
                              type="number"
                              min="1"
                              step="1"
                              onInput=${a=>Q(u=>({...u,form:{...u.form,always_online_interval_seconds:a}}))}
                            />

                            <div className="form-grid">
                              <${pe}
                                label="Proxy protocolo"
                                value=${U.form.proxy_protocol}
                                placeholder="http, https ou socks5"
                                onInput=${a=>Q(u=>({...u,form:{...u.form,proxy_protocol:a}}))}
                              />
                              <${pe}
                                label="Proxy host"
                                value=${U.form.proxy_host}
                                placeholder="proxy.exemplo.com"
                                onInput=${a=>Q(u=>({...u,form:{...u.form,proxy_host:a}}))}
                              />
                              <${pe}
                                label="Proxy porta"
                                value=${U.form.proxy_port}
                                placeholder="8080"
                                onInput=${a=>Q(u=>({...u,form:{...u.form,proxy_port:a}}))}
                              />
                              <${pe}
                                label="Proxy usuário"
                                value=${U.form.proxy_username}
                                placeholder="Opcional"
                                onInput=${a=>Q(u=>({...u,form:{...u.form,proxy_username:a}}))}
                              />
                              <${pe}
                                label="Proxy senha"
                                value=${U.form.proxy_password}
                                type="password"
                                placeholder="Opcional"
                                onInput=${a=>Q(u=>({...u,form:{...u.form,proxy_password:a}}))}
                              />
                            </div>

                            <div className="inline-actions inline-actions--wrap">
                              <button
                                type="submit"
                                className="solid-button"
                                disabled=${U.loading||U.saving}
                              >
                                ${U.saving?"Salvando...":"Salvar configuração"}
                              </button>
                              <button
                                type="button"
                                className="ghost-button"
                                disabled=${U.loading}
                                onClick=${()=>b(E)}
                              >
                                Recarregar
                              </button>
                              <button
                                type="button"
                                className="ghost-button ghost-button--danger"
                                disabled=${U.saving}
                                onClick=${Gr}
                              >
                                Limpar override
                              </button>
                            </div>
                          </form>
                        `:null}

                    ${he==="webhooks"?H`
                          <form className="stack-form" onSubmit=${sr}>
                            <div className="form-grid">
                              <${pe}
                                label="Webhook ID"
                                value=${w.form.id}
                                placeholder="vazio = auto"
                                onInput=${a=>d(u=>({...u,form:{...u.form,id:a}}))}
                              />
                              <${pe}
                                label="Nome"
                                value=${w.form.name}
                                placeholder="Chatwoot principal"
                                onInput=${a=>d(u=>({...u,form:{...u.form,name:a}}))}
                              />
                            </div>

                            <${pe}
                              label="URL base"
                              value=${w.form.url}
                              placeholder="https://seu-app.com/webhooks/whatsapp"
                              onInput=${a=>d(u=>({...u,form:{...u.form,url:a}}))}
                            />
                            <${pe}
                              label="URL absoluta"
                              value=${w.form.url_absolute}
                              placeholder="https://seu-app.com/webhooks/whatsapp/5562..."
                              onInput=${a=>d(u=>({...u,form:{...u.form,url_absolute:a}}))}
                            />

                            <div className="form-grid">
                              <${pe}
                                label="Header auth"
                                value=${w.form.header}
                                placeholder="Authorization"
                                onInput=${a=>d(u=>({...u,form:{...u.form,header:a}}))}
                              />
                              <${pe}
                                label="Token"
                                value=${w.form.token}
                                placeholder="Bearer xxxxx"
                                onInput=${a=>d(u=>({...u,form:{...u.form,token:a}}))}
                              />
                            </div>

                            <${pe}
                              label="Timeout (ms)"
                              value=${w.form.timeout_ms}
                              type="number"
                              min="1"
                              step="100"
                              onInput=${a=>d(u=>({...u,form:{...u.form,timeout_ms:a}}))}
                            />

                            <div className="check-grid">
                              <${Ke}
                                label="Webhook habilitado"
                                checked=${w.form.enabled}
                                onChange=${a=>d(u=>({...u,form:{...u.form,enabled:a}}))}
                              />
                              <${Ke}
                                label="Adicionar sessão à URL"
                                checked=${w.form.append_session}
                                onChange=${a=>d(u=>({...u,form:{...u.form,append_session:a}}))}
                              />
                              <${Ke}
                                label="Enviar incoming"
                                checked=${w.form.send_incoming_messages}
                                onChange=${a=>d(u=>({...u,form:{...u.form,send_incoming_messages:a}}))}
                              />
                              <${Ke}
                                label="Enviar outgoing"
                                checked=${w.form.send_outgoing_messages}
                                onChange=${a=>d(u=>({...u,form:{...u.form,send_outgoing_messages:a}}))}
                              />
                              <${Ke}
                                label="Enviar updates"
                                checked=${w.form.send_update_messages}
                                onChange=${a=>d(u=>({...u,form:{...u.form,send_update_messages:a}}))}
                              />
                              <${Ke}
                                label="Enviar grupos"
                                checked=${w.form.send_group_messages}
                                onChange=${a=>d(u=>({...u,form:{...u.form,send_group_messages:a}}))}
                              />
                              <${Ke}
                                label="Enviar newsletters"
                                checked=${w.form.send_newsletter_messages}
                                onChange=${a=>d(u=>({...u,form:{...u.form,send_newsletter_messages:a}}))}
                              />
                            </div>

                            <div className="inline-actions inline-actions--wrap">
                              <button
                                type="submit"
                                className="solid-button"
                                disabled=${w.loading||w.saving}
                              >
                                ${w.saving?"Salvando...":"Salvar webhook"}
                              </button>
                              <button
                                type="button"
                                className="ghost-button"
                                onClick=${()=>d(a=>({...a,form:rr()}))}
                              >
                                Limpar formulário
                              </button>
                              <button
                                type="button"
                                className="ghost-button"
                                disabled=${w.loading}
                                onClick=${()=>ae(E,{resetForm:!1})}
                              >
                                Recarregar lista
                              </button>
                            </div>
                          </form>

                          <div className="subsection">
                            <div className="panel__header panel__header--space-between">
                              <div>
                                <p className="eyebrow">Cadastrados</p>
                                <h3>Webhooks ativos</h3>
                              </div>
                              <span className="panel__hint">${w.items.length} item(ns)</span>
                            </div>

                            ${w.loading?H`
                                  <div className="empty-state empty-state--compact">
                                    <strong>Carregando webhooks…</strong>
                                  </div>
                                `:w.items.length===0?H`
                                    <div className="empty-state empty-state--compact">
                                      <strong>Nenhum webhook cadastrado</strong>
                                      <p>Preencha o formulário acima para salvar o primeiro callback.</p>
                                    </div>
                                  `:H`
                                    <div className="webhook-list">
                                      ${w.items.map(a=>H`
                                          <${gd}
                                            key=${a.id}
                                            item=${a}
                                            onEdit=${()=>d(u=>({...u,form:ad(a)}))}
                                            onDelete=${()=>Ln(a.id)}
                                          />
                                        `)}
                                    </div>
                                  `}
                          </div>
                        `:null}

                    ${he==="templates"?H`
                          <form className="stack-form" onSubmit=${Xr}>
                            <div className="form-grid">
                              <${pe}
                                label="Template ID"
                                value=${y.form.id}
                                placeholder="vazio = nome__idioma"
                                onInput=${a=>D(u=>({...u,form:{...u.form,id:a}}))}
                              />
                              <${pe}
                                label="Nome Cloud API"
                                value=${y.form.name}
                                placeholder="boleto_disponivel"
                                onInput=${a=>D(u=>({...u,form:{...u.form,name:a}}))}
                              />
                              <${pe}
                                label="Idioma"
                                value=${y.form.language}
                                placeholder="pt_BR"
                                onInput=${a=>D(u=>({...u,form:{...u.form,language:a}}))}
                              />
                              <${pe}
                                label="Categoria"
                                value=${y.form.category}
                                placeholder="utility"
                                onInput=${a=>D(u=>({...u,form:{...u.form,category:a}}))}
                              />
                            </div>

                            <div className="form-grid">
                              <label className="field">
                                <span>Tipo do cabeçalho</span>
                                <select
                                  value=${y.form.header_type}
                                  onChange=${a=>D(u=>({...u,form:{...u.form,header_type:a.target.value}}))}
                                >
                                  <option value="">Sem cabeçalho</option>
                                  <option value="text">Texto</option>
                                  <option value="image">Imagem</option>
                                  <option value="document">Documento</option>
                                </select>
                              </label>
                              <${pe}
                                label="Texto do cabeçalho"
                                value=${y.form.header_text}
                                placeholder="Olá {{nome}}"
                                onInput=${a=>D(u=>({...u,form:{...u.form,header_text:a}}))}
                              />
                            </div>

                            ${y.form.header_type==="image"||y.form.header_type==="document"?H`
                                  <div className="form-grid">
                                    <${pe}
                                      label="Link padrão da mídia"
                                      value=${y.form.header_link}
                                      placeholder="https://..."
                                      onInput=${a=>D(u=>({...u,form:{...u.form,header_link:a}}))}
                                    />
                                    <${pe}
                                      label="Media ID padrão"
                                      value=${y.form.header_id}
                                      placeholder="opcional"
                                      onInput=${a=>D(u=>({...u,form:{...u.form,header_id:a}}))}
                                    />
                                    <${pe}
                                      label="Nome do arquivo"
                                      value=${y.form.header_filename}
                                      placeholder="documento.pdf"
                                      onInput=${a=>D(u=>({...u,form:{...u.form,header_filename:a}}))}
                                    />
                                  </div>
                                `:null}

                            <${oo}
                              label="Corpo"
                              value=${y.form.body}
                              rows="5"
                              placeholder="Olá {{nome}}, seu boleto {{1}} está disponível."
                              onInput=${a=>D(u=>({...u,form:{...u.form,body:a}}))}
                            />

                            <${oo}
                              label="Footer"
                              value=${y.form.footer}
                              rows="2"
                              placeholder="Opcional"
                              onInput=${a=>D(u=>({...u,form:{...u.form,footer:a}}))}
                            />

                            <div className="callout callout--muted">
                              <strong>Variáveis aceitas</strong>
                              <p>
                                Use placeholders como {{'{{1}}'}} ou {{'{{nome}}'}}. No envio, eles podem
                                vir em components.parameters por ordem, com parameter_name, ou no objeto
                                template.variables.
                              </p>
                            </div>

                            <div className="inline-actions inline-actions--wrap">
                              <button
                                type="submit"
                                className="solid-button"
                                disabled=${y.loading||y.saving}
                              >
                                ${y.saving?"Salvando...":"Salvar template"}
                              </button>
                              <button
                                type="button"
                                className="ghost-button"
                                onClick=${()=>D(a=>({...a,form:lr()}))}
                              >
                                Limpar formulário
                              </button>
                              <button
                                type="button"
                                className="ghost-button"
                                disabled=${y.loading}
                                onClick=${()=>Me(E,{resetForm:!1})}
                              >
                                Recarregar lista
                              </button>
                            </div>
                          </form>

                          <div className="subsection">
                            <div className="panel__header panel__header--space-between">
                              <div>
                                <p className="eyebrow">Cadastrados</p>
                                <h3>Templates da sessão</h3>
                              </div>
                              <span className="panel__hint">${y.items.length} item(ns)</span>
                            </div>

                            ${y.loading?H`
                                  <div className="empty-state empty-state--compact">
                                    <strong>Carregando templates…</strong>
                                  </div>
                                `:y.items.length===0?H`
                                    <div className="empty-state empty-state--compact">
                                      <strong>Nenhum template cadastrado</strong>
                                      <p>Preencha o formulário acima para habilitar envios type=template.</p>
                                    </div>
                                  `:H`
                                    <div className="webhook-list">
                                      ${y.items.map(a=>H`
                                          <${yd}
                                            key=${a.id}
                                            item=${a}
                                            onEdit=${()=>D(u=>({...u,form:ud(a)}))}
                                            onDelete=${()=>hn(a.id)}
                                          />
                                        `)}
                                    </div>
                                  `}
                          </div>
                        `:null}
                  `:H`
                    <div className="empty-state">
                      <strong>Nenhuma sessão selecionada</strong>
                      <p>
                        Escolha uma sessão na coluna central para abrir o resumo, a configuração e os
                        webhooks dessa sessão.
                      </p>
                    </div>
                  `}
            </section>
          </div>
        </section>
          `:null}
      </main>

      <div className="toast-stack">
        ${Oe.map(a=>H`
            <div key=${a.id} className=${`toast toast--${a.tone}`}>
              <strong>${a.tone==="error"?"Erro":a.tone==="success"?"OK":"Info"}</strong>
              <span>${a.message}</span>
            </div>
          `)}
      </div>

      <${hd} sessionId=${gt} onClose=${()=>yt("")} />
    </div>
  `}const Ku=document.getElementById("root");Ku&&ed.createRoot(Ku).render(H`<${_d} />`);
