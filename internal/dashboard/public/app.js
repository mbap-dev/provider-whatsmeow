(function(){const _=document.createElement("link").relList;if(_&&_.supports&&_.supports("modulepreload"))return;for(const E of document.querySelectorAll('link[rel="modulepreload"]'))T(E);new MutationObserver(E=>{for(const x of E)if(x.type==="childList")for(const M of x.addedNodes)M.tagName==="LINK"&&M.rel==="modulepreload"&&T(M)}).observe(document,{childList:!0,subtree:!0});function m(E){const x={};return E.integrity&&(x.integrity=E.integrity),E.referrerPolicy&&(x.referrerPolicy=E.referrerPolicy),E.crossOrigin==="use-credentials"?x.credentials="include":E.crossOrigin==="anonymous"?x.credentials="omit":x.credentials="same-origin",x}function T(E){if(E.ep)return;E.ep=!0;const x=m(E);fetch(E.href,x)}})();function Kd(c){return c&&c.__esModule&&Object.prototype.hasOwnProperty.call(c,"default")?c.default:c}var Yi={exports:{}},q={};var Ou;function Yd(){if(Ou)return q;Ou=1;var c=Symbol.for("react.element"),_=Symbol.for("react.portal"),m=Symbol.for("react.fragment"),T=Symbol.for("react.strict_mode"),E=Symbol.for("react.profiler"),x=Symbol.for("react.provider"),M=Symbol.for("react.context"),G=Symbol.for("react.forward_ref"),C=Symbol.for("react.suspense"),ue=Symbol.for("react.memo"),ye=Symbol.for("react.lazy"),le=Symbol.iterator;function Z(f){return f===null||typeof f!="object"?null:(f=le&&f[le]||f["@@iterator"],typeof f=="function"?f:null)}var ut={isMounted:function(){return!1},enqueueForceUpdate:function(){},enqueueReplaceState:function(){},enqueueSetState:function(){}},Oe=Object.assign,fe={};function oe(f,y,D){this.props=f,this.context=y,this.refs=fe,this.updater=D||ut}oe.prototype.isReactComponent={},oe.prototype.setState=function(f,y){if(typeof f!="object"&&typeof f!="function"&&f!=null)throw Error("setState(...): takes an object of state variables to update or a function which returns an object of state variables.");this.updater.enqueueSetState(this,f,y,"setState")},oe.prototype.forceUpdate=function(f){this.updater.enqueueForceUpdate(this,f,"forceUpdate")};function vt(){}vt.prototype=oe.prototype;function tt(f,y,D){this.props=f,this.context=y,this.refs=fe,this.updater=D||ut}var nt=tt.prototype=new vt;nt.constructor=tt,Oe(nt,oe.prototype),nt.isPureReactComponent=!0;var Se=Array.isArray,rt=Object.prototype.hasOwnProperty,Fe={current:null},Me={key:!0,ref:!0,__self:!0,__source:!0};function Be(f,y,D){var Y,z={},X=null,b=null;if(y!=null)for(Y in y.ref!==void 0&&(b=y.ref),y.key!==void 0&&(X=""+y.key),y)rt.call(y,Y)&&!Me.hasOwnProperty(Y)&&(z[Y]=y[Y]);var ee=arguments.length-2;if(ee===1)z.children=D;else if(1<ee){for(var ce=Array(ee),je=0;je<ee;je++)ce[je]=arguments[je+2];z.children=ce}if(f&&f.defaultProps)for(Y in ee=f.defaultProps,ee)z[Y]===void 0&&(z[Y]=ee[Y]);return{$$typeof:c,type:f,key:X,ref:b,props:z,_owner:Fe.current}}function xt(f,y){return{$$typeof:c,type:f.type,key:y,ref:f.ref,props:f.props,_owner:f._owner}}function gt(f){return typeof f=="object"&&f!==null&&f.$$typeof===c}function yt(f){var y={"=":"=0",":":"=2"};return"$"+f.replace(/[=:]/g,function(D){return y[D]})}var J=/\/+/g;function me(f,y){return typeof f=="object"&&f!==null&&f.key!=null?yt(""+f.key):y.toString(36)}function Pe(f,y,D,Y,z){var X=typeof f;(X==="undefined"||X==="boolean")&&(f=null);var b=!1;if(f===null)b=!0;else switch(X){case"string":case"number":b=!0;break;case"object":switch(f.$$typeof){case c:case _:b=!0}}if(b)return b=f,z=z(b),f=Y===""?"."+me(b,0):Y,Se(z)?(D="",f!=null&&(D=f.replace(J,"$&/")+"/"),Pe(z,y,D,"",function(je){return je})):z!=null&&(gt(z)&&(z=xt(z,D+(!z.key||b&&b.key===z.key?"":(""+z.key).replace(J,"$&/")+"/")+f)),y.push(z)),1;if(b=0,Y=Y===""?".":Y+":",Se(f))for(var ee=0;ee<f.length;ee++){X=f[ee];var ce=Y+me(X,ee);b+=Pe(X,y,D,ce,z)}else if(ce=Z(f),typeof ce=="function")for(f=ce.call(f),ee=0;!(X=f.next()).done;)X=X.value,ce=Y+me(X,ee++),b+=Pe(X,y,D,ce,z);else if(X==="object")throw y=String(f),Error("Objects are not valid as a React child (found: "+(y==="[object Object]"?"object with keys {"+Object.keys(f).join(", ")+"}":y)+"). If you meant to render a collection of children, use an array instead.");return b}function We(f,y,D){if(f==null)return f;var Y=[],z=0;return Pe(f,Y,"","",function(X){return y.call(D,X,z++)}),Y}function j(f){if(f._status===-1){var y=f._result;y=y(),y.then(function(D){(f._status===0||f._status===-1)&&(f._status=1,f._result=D)},function(D){(f._status===0||f._status===-1)&&(f._status=2,f._result=D)}),f._status===-1&&(f._status=0,f._result=y)}if(f._status===1)return f._result.default;throw f._result}var Q={current:null},P={transition:null},U={ReactCurrentDispatcher:Q,ReactCurrentBatchConfig:P,ReactCurrentOwner:Fe};function $(){throw Error("act(...) is not supported in production builds of React.")}return q.Children={map:We,forEach:function(f,y,D){We(f,function(){y.apply(this,arguments)},D)},count:function(f){var y=0;return We(f,function(){y++}),y},toArray:function(f){return We(f,function(y){return y})||[]},only:function(f){if(!gt(f))throw Error("React.Children.only expected to receive a single React element child.");return f}},q.Component=oe,q.Fragment=m,q.Profiler=E,q.PureComponent=tt,q.StrictMode=T,q.Suspense=C,q.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED=U,q.act=$,q.cloneElement=function(f,y,D){if(f==null)throw Error("React.cloneElement(...): The argument must be a React element, but you passed "+f+".");var Y=Oe({},f.props),z=f.key,X=f.ref,b=f._owner;if(y!=null){if(y.ref!==void 0&&(X=y.ref,b=Fe.current),y.key!==void 0&&(z=""+y.key),f.type&&f.type.defaultProps)var ee=f.type.defaultProps;for(ce in y)rt.call(y,ce)&&!Me.hasOwnProperty(ce)&&(Y[ce]=y[ce]===void 0&&ee!==void 0?ee[ce]:y[ce])}var ce=arguments.length-2;if(ce===1)Y.children=D;else if(1<ce){ee=Array(ce);for(var je=0;je<ce;je++)ee[je]=arguments[je+2];Y.children=ee}return{$$typeof:c,type:f.type,key:z,ref:X,props:Y,_owner:b}},q.createContext=function(f){return f={$$typeof:M,_currentValue:f,_currentValue2:f,_threadCount:0,Provider:null,Consumer:null,_defaultValue:null,_globalName:null},f.Provider={$$typeof:x,_context:f},f.Consumer=f},q.createElement=Be,q.createFactory=function(f){var y=Be.bind(null,f);return y.type=f,y},q.createRef=function(){return{current:null}},q.forwardRef=function(f){return{$$typeof:G,render:f}},q.isValidElement=gt,q.lazy=function(f){return{$$typeof:ye,_payload:{_status:-1,_result:f},_init:j}},q.memo=function(f,y){return{$$typeof:ue,type:f,compare:y===void 0?null:y}},q.startTransition=function(f){var y=P.transition;P.transition={};try{f()}finally{P.transition=y}},q.unstable_act=$,q.useCallback=function(f,y){return Q.current.useCallback(f,y)},q.useContext=function(f){return Q.current.useContext(f)},q.useDebugValue=function(){},q.useDeferredValue=function(f){return Q.current.useDeferredValue(f)},q.useEffect=function(f,y){return Q.current.useEffect(f,y)},q.useId=function(){return Q.current.useId()},q.useImperativeHandle=function(f,y,D){return Q.current.useImperativeHandle(f,y,D)},q.useInsertionEffect=function(f,y){return Q.current.useInsertionEffect(f,y)},q.useLayoutEffect=function(f,y){return Q.current.useLayoutEffect(f,y)},q.useMemo=function(f,y){return Q.current.useMemo(f,y)},q.useReducer=function(f,y,D){return Q.current.useReducer(f,y,D)},q.useRef=function(f){return Q.current.useRef(f)},q.useState=function(f){return Q.current.useState(f)},q.useSyncExternalStore=function(f,y,D){return Q.current.useSyncExternalStore(f,y,D)},q.useTransition=function(){return Q.current.useTransition()},q.version="18.3.1",q}var Fu;function Gu(){return Fu||(Fu=1,Yi.exports=Yd()),Yi.exports}var ie=Gu();const Gd=Kd(ie);var ro={},Gi={exports:{}},et={},Xi={exports:{}},qi={};var Mu;function Xd(){return Mu||(Mu=1,(function(c){function _(P,U){var $=P.length;P.push(U);e:for(;0<$;){var f=$-1>>>1,y=P[f];if(0<E(y,U))P[f]=U,P[$]=y,$=f;else break e}}function m(P){return P.length===0?null:P[0]}function T(P){if(P.length===0)return null;var U=P[0],$=P.pop();if($!==U){P[0]=$;e:for(var f=0,y=P.length,D=y>>>1;f<D;){var Y=2*(f+1)-1,z=P[Y],X=Y+1,b=P[X];if(0>E(z,$))X<y&&0>E(b,z)?(P[f]=b,P[X]=$,f=X):(P[f]=z,P[Y]=$,f=Y);else if(X<y&&0>E(b,$))P[f]=b,P[X]=$,f=X;else break e}}return U}function E(P,U){var $=P.sortIndex-U.sortIndex;return $!==0?$:P.id-U.id}if(typeof performance=="object"&&typeof performance.now=="function"){var x=performance;c.unstable_now=function(){return x.now()}}else{var M=Date,G=M.now();c.unstable_now=function(){return M.now()-G}}var C=[],ue=[],ye=1,le=null,Z=3,ut=!1,Oe=!1,fe=!1,oe=typeof setTimeout=="function"?setTimeout:null,vt=typeof clearTimeout=="function"?clearTimeout:null,tt=typeof setImmediate<"u"?setImmediate:null;typeof navigator<"u"&&navigator.scheduling!==void 0&&navigator.scheduling.isInputPending!==void 0&&navigator.scheduling.isInputPending.bind(navigator.scheduling);function nt(P){for(var U=m(ue);U!==null;){if(U.callback===null)T(ue);else if(U.startTime<=P)T(ue),U.sortIndex=U.expirationTime,_(C,U);else break;U=m(ue)}}function Se(P){if(fe=!1,nt(P),!Oe)if(m(C)!==null)Oe=!0,j(rt);else{var U=m(ue);U!==null&&Q(Se,U.startTime-P)}}function rt(P,U){Oe=!1,fe&&(fe=!1,vt(Be),Be=-1),ut=!0;var $=Z;try{for(nt(U),le=m(C);le!==null&&(!(le.expirationTime>U)||P&&!yt());){var f=le.callback;if(typeof f=="function"){le.callback=null,Z=le.priorityLevel;var y=f(le.expirationTime<=U);U=c.unstable_now(),typeof y=="function"?le.callback=y:le===m(C)&&T(C),nt(U)}else T(C);le=m(C)}if(le!==null)var D=!0;else{var Y=m(ue);Y!==null&&Q(Se,Y.startTime-U),D=!1}return D}finally{le=null,Z=$,ut=!1}}var Fe=!1,Me=null,Be=-1,xt=5,gt=-1;function yt(){return!(c.unstable_now()-gt<xt)}function J(){if(Me!==null){var P=c.unstable_now();gt=P;var U=!0;try{U=Me(!0,P)}finally{U?me():(Fe=!1,Me=null)}}else Fe=!1}var me;if(typeof tt=="function")me=function(){tt(J)};else if(typeof MessageChannel<"u"){var Pe=new MessageChannel,We=Pe.port2;Pe.port1.onmessage=J,me=function(){We.postMessage(null)}}else me=function(){oe(J,0)};function j(P){Me=P,Fe||(Fe=!0,me())}function Q(P,U){Be=oe(function(){P(c.unstable_now())},U)}c.unstable_IdlePriority=5,c.unstable_ImmediatePriority=1,c.unstable_LowPriority=4,c.unstable_NormalPriority=3,c.unstable_Profiling=null,c.unstable_UserBlockingPriority=2,c.unstable_cancelCallback=function(P){P.callback=null},c.unstable_continueExecution=function(){Oe||ut||(Oe=!0,j(rt))},c.unstable_forceFrameRate=function(P){0>P||125<P?console.error("forceFrameRate takes a positive int between 0 and 125, forcing frame rates higher than 125 fps is not supported"):xt=0<P?Math.floor(1e3/P):5},c.unstable_getCurrentPriorityLevel=function(){return Z},c.unstable_getFirstCallbackNode=function(){return m(C)},c.unstable_next=function(P){switch(Z){case 1:case 2:case 3:var U=3;break;default:U=Z}var $=Z;Z=U;try{return P()}finally{Z=$}},c.unstable_pauseExecution=function(){},c.unstable_requestPaint=function(){},c.unstable_runWithPriority=function(P,U){switch(P){case 1:case 2:case 3:case 4:case 5:break;default:P=3}var $=Z;Z=P;try{return U()}finally{Z=$}},c.unstable_scheduleCallback=function(P,U,$){var f=c.unstable_now();switch(typeof $=="object"&&$!==null?($=$.delay,$=typeof $=="number"&&0<$?f+$:f):$=f,P){case 1:var y=-1;break;case 2:y=250;break;case 5:y=1073741823;break;case 4:y=1e4;break;default:y=5e3}return y=$+y,P={id:ye++,callback:U,priorityLevel:P,startTime:$,expirationTime:y,sortIndex:-1},$>f?(P.sortIndex=$,_(ue,P),m(C)===null&&P===m(ue)&&(fe?(vt(Be),Be=-1):fe=!0,Q(Se,$-f))):(P.sortIndex=y,_(C,P),Oe||ut||(Oe=!0,j(rt))),P},c.unstable_shouldYield=yt,c.unstable_wrapCallback=function(P){var U=Z;return function(){var $=Z;Z=U;try{return P.apply(this,arguments)}finally{Z=$}}}})(qi)),qi}var ju;function qd(){return ju||(ju=1,Xi.exports=Xd()),Xi.exports}var Au;function Zd(){if(Au)return et;Au=1;var c=Gu(),_=qd();function m(e){for(var t="https://reactjs.org/docs/error-decoder.html?invariant="+e,n=1;n<arguments.length;n++)t+="&args[]="+encodeURIComponent(arguments[n]);return"Minified React error #"+e+"; visit "+t+" for the full message or use the non-minified dev environment for full errors and additional helpful warnings."}var T=new Set,E={};function x(e,t){M(e,t),M(e+"Capture",t)}function M(e,t){for(E[e]=t,e=0;e<t.length;e++)T.add(t[e])}var G=!(typeof window>"u"||typeof window.document>"u"||typeof window.document.createElement>"u"),C=Object.prototype.hasOwnProperty,ue=/^[:A-Z_a-z\u00C0-\u00D6\u00D8-\u00F6\u00F8-\u02FF\u0370-\u037D\u037F-\u1FFF\u200C-\u200D\u2070-\u218F\u2C00-\u2FEF\u3001-\uD7FF\uF900-\uFDCF\uFDF0-\uFFFD][:A-Z_a-z\u00C0-\u00D6\u00D8-\u00F6\u00F8-\u02FF\u0370-\u037D\u037F-\u1FFF\u200C-\u200D\u2070-\u218F\u2C00-\u2FEF\u3001-\uD7FF\uF900-\uFDCF\uFDF0-\uFFFD\-.0-9\u00B7\u0300-\u036F\u203F-\u2040]*$/,ye={},le={};function Z(e){return C.call(le,e)?!0:C.call(ye,e)?!1:ue.test(e)?le[e]=!0:(ye[e]=!0,!1)}function ut(e,t,n,r){if(n!==null&&n.type===0)return!1;switch(typeof t){case"function":case"symbol":return!0;case"boolean":return r?!1:n!==null?!n.acceptsBooleans:(e=e.toLowerCase().slice(0,5),e!=="data-"&&e!=="aria-");default:return!1}}function Oe(e,t,n,r){if(t===null||typeof t>"u"||ut(e,t,n,r))return!0;if(r)return!1;if(n!==null)switch(n.type){case 3:return!t;case 4:return t===!1;case 5:return isNaN(t);case 6:return isNaN(t)||1>t}return!1}function fe(e,t,n,r,l,o,i){this.acceptsBooleans=t===2||t===3||t===4,this.attributeName=r,this.attributeNamespace=l,this.mustUseProperty=n,this.propertyName=e,this.type=t,this.sanitizeURL=o,this.removeEmptyString=i}var oe={};"children dangerouslySetInnerHTML defaultValue defaultChecked innerHTML suppressContentEditableWarning suppressHydrationWarning style".split(" ").forEach(function(e){oe[e]=new fe(e,0,!1,e,null,!1,!1)}),[["acceptCharset","accept-charset"],["className","class"],["htmlFor","for"],["httpEquiv","http-equiv"]].forEach(function(e){var t=e[0];oe[t]=new fe(t,1,!1,e[1],null,!1,!1)}),["contentEditable","draggable","spellCheck","value"].forEach(function(e){oe[e]=new fe(e,2,!1,e.toLowerCase(),null,!1,!1)}),["autoReverse","externalResourcesRequired","focusable","preserveAlpha"].forEach(function(e){oe[e]=new fe(e,2,!1,e,null,!1,!1)}),"allowFullScreen async autoFocus autoPlay controls default defer disabled disablePictureInPicture disableRemotePlayback formNoValidate hidden loop noModule noValidate open playsInline readOnly required reversed scoped seamless itemScope".split(" ").forEach(function(e){oe[e]=new fe(e,3,!1,e.toLowerCase(),null,!1,!1)}),["checked","multiple","muted","selected"].forEach(function(e){oe[e]=new fe(e,3,!0,e,null,!1,!1)}),["capture","download"].forEach(function(e){oe[e]=new fe(e,4,!1,e,null,!1,!1)}),["cols","rows","size","span"].forEach(function(e){oe[e]=new fe(e,6,!1,e,null,!1,!1)}),["rowSpan","start"].forEach(function(e){oe[e]=new fe(e,5,!1,e.toLowerCase(),null,!1,!1)});var vt=/[\-:]([a-z])/g;function tt(e){return e[1].toUpperCase()}"accent-height alignment-baseline arabic-form baseline-shift cap-height clip-path clip-rule color-interpolation color-interpolation-filters color-profile color-rendering dominant-baseline enable-background fill-opacity fill-rule flood-color flood-opacity font-family font-size font-size-adjust font-stretch font-style font-variant font-weight glyph-name glyph-orientation-horizontal glyph-orientation-vertical horiz-adv-x horiz-origin-x image-rendering letter-spacing lighting-color marker-end marker-mid marker-start overline-position overline-thickness paint-order panose-1 pointer-events rendering-intent shape-rendering stop-color stop-opacity strikethrough-position strikethrough-thickness stroke-dasharray stroke-dashoffset stroke-linecap stroke-linejoin stroke-miterlimit stroke-opacity stroke-width text-anchor text-decoration text-rendering underline-position underline-thickness unicode-bidi unicode-range units-per-em v-alphabetic v-hanging v-ideographic v-mathematical vector-effect vert-adv-y vert-origin-x vert-origin-y word-spacing writing-mode xmlns:xlink x-height".split(" ").forEach(function(e){var t=e.replace(vt,tt);oe[t]=new fe(t,1,!1,e,null,!1,!1)}),"xlink:actuate xlink:arcrole xlink:role xlink:show xlink:title xlink:type".split(" ").forEach(function(e){var t=e.replace(vt,tt);oe[t]=new fe(t,1,!1,e,"http://www.w3.org/1999/xlink",!1,!1)}),["xml:base","xml:lang","xml:space"].forEach(function(e){var t=e.replace(vt,tt);oe[t]=new fe(t,1,!1,e,"http://www.w3.org/XML/1998/namespace",!1,!1)}),["tabIndex","crossOrigin"].forEach(function(e){oe[e]=new fe(e,1,!1,e.toLowerCase(),null,!1,!1)}),oe.xlinkHref=new fe("xlinkHref",1,!1,"xlink:href","http://www.w3.org/1999/xlink",!0,!1),["src","href","action","formAction"].forEach(function(e){oe[e]=new fe(e,1,!1,e.toLowerCase(),null,!0,!0)});function nt(e,t,n,r){var l=oe.hasOwnProperty(t)?oe[t]:null;(l!==null?l.type!==0:r||!(2<t.length)||t[0]!=="o"&&t[0]!=="O"||t[1]!=="n"&&t[1]!=="N")&&(Oe(t,n,l,r)&&(n=null),r||l===null?Z(t)&&(n===null?e.removeAttribute(t):e.setAttribute(t,""+n)):l.mustUseProperty?e[l.propertyName]=n===null?l.type===3?!1:"":n:(t=l.attributeName,r=l.attributeNamespace,n===null?e.removeAttribute(t):(l=l.type,n=l===3||l===4&&n===!0?"":""+n,r?e.setAttributeNS(r,t,n):e.setAttribute(t,n))))}var Se=c.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED,rt=Symbol.for("react.element"),Fe=Symbol.for("react.portal"),Me=Symbol.for("react.fragment"),Be=Symbol.for("react.strict_mode"),xt=Symbol.for("react.profiler"),gt=Symbol.for("react.provider"),yt=Symbol.for("react.context"),J=Symbol.for("react.forward_ref"),me=Symbol.for("react.suspense"),Pe=Symbol.for("react.suspense_list"),We=Symbol.for("react.memo"),j=Symbol.for("react.lazy"),Q=Symbol.for("react.offscreen"),P=Symbol.iterator;function U(e){return e===null||typeof e!="object"?null:(e=P&&e[P]||e["@@iterator"],typeof e=="function"?e:null)}var $=Object.assign,f;function y(e){if(f===void 0)try{throw Error()}catch(n){var t=n.stack.trim().match(/\n( *(at )?)/);f=t&&t[1]||""}return`
`+f+e}var D=!1;function Y(e,t){if(!e||D)return"";D=!0;var n=Error.prepareStackTrace;Error.prepareStackTrace=void 0;try{if(t)if(t=function(){throw Error()},Object.defineProperty(t.prototype,"props",{set:function(){throw Error()}}),typeof Reflect=="object"&&Reflect.construct){try{Reflect.construct(t,[])}catch(g){var r=g}Reflect.construct(e,[],t)}else{try{t.call()}catch(g){r=g}e.call(t.prototype)}else{try{throw Error()}catch(g){r=g}e()}}catch(g){if(g&&r&&typeof g.stack=="string"){for(var l=g.stack.split(`
`),o=r.stack.split(`
`),i=l.length-1,s=o.length-1;1<=i&&0<=s&&l[i]!==o[s];)s--;for(;1<=i&&0<=s;i--,s--)if(l[i]!==o[s]){if(i!==1||s!==1)do if(i--,s--,0>s||l[i]!==o[s]){var d=`
`+l[i].replace(" at new "," at ");return e.displayName&&d.includes("<anonymous>")&&(d=d.replace("<anonymous>",e.displayName)),d}while(1<=i&&0<=s);break}}}finally{D=!1,Error.prepareStackTrace=n}return(e=e?e.displayName||e.name:"")?y(e):""}function z(e){switch(e.tag){case 5:return y(e.type);case 16:return y("Lazy");case 13:return y("Suspense");case 19:return y("SuspenseList");case 0:case 2:case 15:return e=Y(e.type,!1),e;case 11:return e=Y(e.type.render,!1),e;case 1:return e=Y(e.type,!0),e;default:return""}}function X(e){if(e==null)return null;if(typeof e=="function")return e.displayName||e.name||null;if(typeof e=="string")return e;switch(e){case Me:return"Fragment";case Fe:return"Portal";case xt:return"Profiler";case Be:return"StrictMode";case me:return"Suspense";case Pe:return"SuspenseList"}if(typeof e=="object")switch(e.$$typeof){case yt:return(e.displayName||"Context")+".Consumer";case gt:return(e._context.displayName||"Context")+".Provider";case J:var t=e.render;return e=e.displayName,e||(e=t.displayName||t.name||"",e=e!==""?"ForwardRef("+e+")":"ForwardRef"),e;case We:return t=e.displayName||null,t!==null?t:X(e.type)||"Memo";case j:t=e._payload,e=e._init;try{return X(e(t))}catch{}}return null}function b(e){var t=e.type;switch(e.tag){case 24:return"Cache";case 9:return(t.displayName||"Context")+".Consumer";case 10:return(t._context.displayName||"Context")+".Provider";case 18:return"DehydratedFragment";case 11:return e=t.render,e=e.displayName||e.name||"",t.displayName||(e!==""?"ForwardRef("+e+")":"ForwardRef");case 7:return"Fragment";case 5:return t;case 4:return"Portal";case 3:return"Root";case 6:return"Text";case 16:return X(t);case 8:return t===Be?"StrictMode":"Mode";case 22:return"Offscreen";case 12:return"Profiler";case 21:return"Scope";case 13:return"Suspense";case 19:return"SuspenseList";case 25:return"TracingMarker";case 1:case 0:case 17:case 2:case 14:case 15:if(typeof t=="function")return t.displayName||t.name||null;if(typeof t=="string")return t}return null}function ee(e){switch(typeof e){case"boolean":case"number":case"string":case"undefined":return e;case"object":return e;default:return""}}function ce(e){var t=e.type;return(e=e.nodeName)&&e.toLowerCase()==="input"&&(t==="checkbox"||t==="radio")}function je(e){var t=ce(e)?"checked":"value",n=Object.getOwnPropertyDescriptor(e.constructor.prototype,t),r=""+e[t];if(!e.hasOwnProperty(t)&&typeof n<"u"&&typeof n.get=="function"&&typeof n.set=="function"){var l=n.get,o=n.set;return Object.defineProperty(e,t,{configurable:!0,get:function(){return l.call(this)},set:function(i){r=""+i,o.call(this,i)}}),Object.defineProperty(e,t,{enumerable:n.enumerable}),{getValue:function(){return r},setValue:function(i){r=""+i},stopTracking:function(){e._valueTracker=null,delete e[t]}}}}function se(e){e._valueTracker||(e._valueTracker=je(e))}function pn(e){if(!e)return!1;var t=e._valueTracker;if(!t)return!0;var n=t.getValue(),r="";return e&&(r=ce(e)?e.checked?"true":"false":e.value),e=r,e!==n?(t.setValue(e),!0):!1}function mn(e){if(e=e||(typeof document<"u"?document:void 0),typeof e>"u")return null;try{return e.activeElement||e.body}catch{return e.body}}function Qt(e,t){var n=t.checked;return $({},t,{defaultChecked:void 0,defaultValue:void 0,value:void 0,checked:n??e._wrapperState.initialChecked})}function Vr(e,t){var n=t.defaultValue==null?"":t.defaultValue,r=t.checked!=null?t.checked:t.defaultChecked;n=ee(t.value!=null?t.value:n),e._wrapperState={initialChecked:r,initialValue:n,controlled:t.type==="checkbox"||t.type==="radio"?t.checked!=null:t.value!=null}}function or(e,t){t=t.checked,t!=null&&nt(e,"checked",t,!1)}function hn(e,t){or(e,t);var n=ee(t.value),r=t.type;if(n!=null)r==="number"?(n===0&&e.value===""||e.value!=n)&&(e.value=""+n):e.value!==""+n&&(e.value=""+n);else if(r==="submit"||r==="reset"){e.removeAttribute("value");return}t.hasOwnProperty("value")?ir(e,t.type,n):t.hasOwnProperty("defaultValue")&&ir(e,t.type,ee(t.defaultValue)),t.checked==null&&t.defaultChecked!=null&&(e.defaultChecked=!!t.defaultChecked)}function Ln(e,t,n){if(t.hasOwnProperty("value")||t.hasOwnProperty("defaultValue")){var r=t.type;if(!(r!=="submit"&&r!=="reset"||t.value!==void 0&&t.value!==null))return;t=""+e._wrapperState.initialValue,n||t===e.value||(e.value=t),e.defaultValue=t}n=e.name,n!==""&&(e.name=""),e.defaultChecked=!!e._wrapperState.initialChecked,n!==""&&(e.name=n)}function ir(e,t,n){(t!=="number"||mn(e.ownerDocument)!==e)&&(n==null?e.defaultValue=""+e._wrapperState.initialValue:e.defaultValue!==""+n&&(e.defaultValue=""+n))}var vn=Array.isArray;function Vt(e,t,n,r){if(e=e.options,t){t={};for(var l=0;l<n.length;l++)t["$"+n[l]]=!0;for(n=0;n<e.length;n++)l=t.hasOwnProperty("$"+e[n].value),e[n].selected!==l&&(e[n].selected=l),l&&r&&(e[n].defaultSelected=!0)}else{for(n=""+ee(n),t=null,l=0;l<e.length;l++){if(e[l].value===n){e[l].selected=!0,r&&(e[l].defaultSelected=!0);return}t!==null||e[l].disabled||(t=e[l])}t!==null&&(t.selected=!0)}}function It(e,t){if(t.dangerouslySetInnerHTML!=null)throw Error(m(91));return $({},t,{value:void 0,defaultValue:void 0,children:""+e._wrapperState.initialValue})}function Hr(e,t){var n=t.value;if(n==null){if(n=t.children,t=t.defaultValue,n!=null){if(t!=null)throw Error(m(92));if(vn(n)){if(1<n.length)throw Error(m(93));n=n[0]}t=n}t==null&&(t=""),n=t}e._wrapperState={initialValue:ee(n)}}function Kr(e,t){var n=ee(t.value),r=ee(t.defaultValue);n!=null&&(n=""+n,n!==e.value&&(e.value=n),t.defaultValue==null&&e.defaultValue!==n&&(e.defaultValue=n)),r!=null&&(e.defaultValue=""+r)}function Yr(e){var t=e.textContent;t===e._wrapperState.initialValue&&t!==""&&t!==null&&(e.value=t)}function Gr(e){switch(e){case"svg":return"http://www.w3.org/2000/svg";case"math":return"http://www.w3.org/1998/Math/MathML";default:return"http://www.w3.org/1999/xhtml"}}function sr(e,t){return e==null||e==="http://www.w3.org/1999/xhtml"?Gr(t):e==="http://www.w3.org/2000/svg"&&t==="foreignObject"?"http://www.w3.org/1999/xhtml":e}var In,Xr=(function(e){return typeof MSApp<"u"&&MSApp.execUnsafeLocalFunction?function(t,n,r,l){MSApp.execUnsafeLocalFunction(function(){return e(t,n,r,l)})}:e})(function(e,t){if(e.namespaceURI!=="http://www.w3.org/2000/svg"||"innerHTML"in e)e.innerHTML=t;else{for(In=In||document.createElement("div"),In.innerHTML="<svg>"+t.valueOf().toString()+"</svg>",t=In.firstChild;e.firstChild;)e.removeChild(e.firstChild);for(;t.firstChild;)e.appendChild(t.firstChild)}});function gn(e,t){if(t){var n=e.firstChild;if(n&&n===e.lastChild&&n.nodeType===3){n.nodeValue=t;return}}e.textContent=t}var yn={animationIterationCount:!0,aspectRatio:!0,borderImageOutset:!0,borderImageSlice:!0,borderImageWidth:!0,boxFlex:!0,boxFlexGroup:!0,boxOrdinalGroup:!0,columnCount:!0,columns:!0,flex:!0,flexGrow:!0,flexPositive:!0,flexShrink:!0,flexNegative:!0,flexOrder:!0,gridArea:!0,gridRow:!0,gridRowEnd:!0,gridRowSpan:!0,gridRowStart:!0,gridColumn:!0,gridColumnEnd:!0,gridColumnSpan:!0,gridColumnStart:!0,fontWeight:!0,lineClamp:!0,lineHeight:!0,opacity:!0,order:!0,orphans:!0,tabSize:!0,widows:!0,zIndex:!0,zoom:!0,fillOpacity:!0,floodOpacity:!0,stopOpacity:!0,strokeDasharray:!0,strokeDashoffset:!0,strokeMiterlimit:!0,strokeOpacity:!0,strokeWidth:!0},a=["Webkit","ms","Moz","O"];Object.keys(yn).forEach(function(e){a.forEach(function(t){t=t+e.charAt(0).toUpperCase()+e.substring(1),yn[t]=yn[e]})});function u(e,t,n){return t==null||typeof t=="boolean"||t===""?"":n||typeof t!="number"||t===0||yn.hasOwnProperty(e)&&yn[e]?(""+t).trim():t+"px"}function L(e,t){e=e.style;for(var n in t)if(t.hasOwnProperty(n)){var r=n.indexOf("--")===0,l=u(n,t[n],r);n==="float"&&(n="cssFloat"),r?e.setProperty(n,l):e[n]=l}}var V=$({menuitem:!0},{area:!0,base:!0,br:!0,col:!0,embed:!0,hr:!0,img:!0,input:!0,keygen:!0,link:!0,meta:!0,param:!0,source:!0,track:!0,wbr:!0});function de(e,t){if(t){if(V[e]&&(t.children!=null||t.dangerouslySetInnerHTML!=null))throw Error(m(137,e));if(t.dangerouslySetInnerHTML!=null){if(t.children!=null)throw Error(m(60));if(typeof t.dangerouslySetInnerHTML!="object"||!("__html"in t.dangerouslySetInnerHTML))throw Error(m(61))}if(t.style!=null&&typeof t.style!="object")throw Error(m(62))}}function _t(e,t){if(e.indexOf("-")===-1)return typeof t.is=="string";switch(e){case"annotation-xml":case"color-profile":case"font-face":case"font-face-src":case"font-face-uri":case"font-face-format":case"font-face-name":case"missing-glyph":return!1;default:return!0}}var io=null;function so(e){return e=e.target||e.srcElement||window,e.correspondingUseElement&&(e=e.correspondingUseElement),e.nodeType===3?e.parentNode:e}var ao=null,Dn=null,On=null;function Ji(e){if(e=Tr(e)){if(typeof ao!="function")throw Error(m(280));var t=e.stateNode;t&&(t=yl(t),ao(e.stateNode,e.type,t))}}function bi(e){Dn?On?On.push(e):On=[e]:Dn=e}function es(){if(Dn){var e=Dn,t=On;if(On=Dn=null,Ji(e),t)for(e=0;e<t.length;e++)Ji(t[e])}}function ts(e,t){return e(t)}function ns(){}var uo=!1;function rs(e,t,n){if(uo)return e(t,n);uo=!0;try{return ts(e,t,n)}finally{uo=!1,(Dn!==null||On!==null)&&(ns(),es())}}function ar(e,t){var n=e.stateNode;if(n===null)return null;var r=yl(n);if(r===null)return null;n=r[t];e:switch(t){case"onClick":case"onClickCapture":case"onDoubleClick":case"onDoubleClickCapture":case"onMouseDown":case"onMouseDownCapture":case"onMouseMove":case"onMouseMoveCapture":case"onMouseUp":case"onMouseUpCapture":case"onMouseEnter":(r=!r.disabled)||(e=e.type,r=!(e==="button"||e==="input"||e==="select"||e==="textarea")),e=!r;break e;default:e=!1}if(e)return null;if(n&&typeof n!="function")throw Error(m(231,t,typeof n));return n}var co=!1;if(G)try{var ur={};Object.defineProperty(ur,"passive",{get:function(){co=!0}}),window.addEventListener("test",ur,ur),window.removeEventListener("test",ur,ur)}catch{co=!1}function Ju(e,t,n,r,l,o,i,s,d){var g=Array.prototype.slice.call(arguments,3);try{t.apply(n,g)}catch(k){this.onError(k)}}var cr=!1,qr=null,Zr=!1,fo=null,bu={onError:function(e){cr=!0,qr=e}};function ec(e,t,n,r,l,o,i,s,d){cr=!1,qr=null,Ju.apply(bu,arguments)}function tc(e,t,n,r,l,o,i,s,d){if(ec.apply(this,arguments),cr){if(cr){var g=qr;cr=!1,qr=null}else throw Error(m(198));Zr||(Zr=!0,fo=g)}}function _n(e){var t=e,n=e;if(e.alternate)for(;t.return;)t=t.return;else{e=t;do t=e,(t.flags&4098)!==0&&(n=t.return),e=t.return;while(e)}return t.tag===3?n:null}function ls(e){if(e.tag===13){var t=e.memoizedState;if(t===null&&(e=e.alternate,e!==null&&(t=e.memoizedState)),t!==null)return t.dehydrated}return null}function os(e){if(_n(e)!==e)throw Error(m(188))}function nc(e){var t=e.alternate;if(!t){if(t=_n(e),t===null)throw Error(m(188));return t!==e?null:e}for(var n=e,r=t;;){var l=n.return;if(l===null)break;var o=l.alternate;if(o===null){if(r=l.return,r!==null){n=r;continue}break}if(l.child===o.child){for(o=l.child;o;){if(o===n)return os(l),e;if(o===r)return os(l),t;o=o.sibling}throw Error(m(188))}if(n.return!==r.return)n=l,r=o;else{for(var i=!1,s=l.child;s;){if(s===n){i=!0,n=l,r=o;break}if(s===r){i=!0,r=l,n=o;break}s=s.sibling}if(!i){for(s=o.child;s;){if(s===n){i=!0,n=o,r=l;break}if(s===r){i=!0,r=o,n=l;break}s=s.sibling}if(!i)throw Error(m(189))}}if(n.alternate!==r)throw Error(m(190))}if(n.tag!==3)throw Error(m(188));return n.stateNode.current===n?e:t}function is(e){return e=nc(e),e!==null?ss(e):null}function ss(e){if(e.tag===5||e.tag===6)return e;for(e=e.child;e!==null;){var t=ss(e);if(t!==null)return t;e=e.sibling}return null}var as=_.unstable_scheduleCallback,us=_.unstable_cancelCallback,rc=_.unstable_shouldYield,lc=_.unstable_requestPaint,$e=_.unstable_now,oc=_.unstable_getCurrentPriorityLevel,po=_.unstable_ImmediatePriority,cs=_.unstable_UserBlockingPriority,Jr=_.unstable_NormalPriority,ic=_.unstable_LowPriority,ds=_.unstable_IdlePriority,br=null,Pt=null;function sc(e){if(Pt&&typeof Pt.onCommitFiberRoot=="function")try{Pt.onCommitFiberRoot(br,e,void 0,(e.current.flags&128)===128)}catch{}}var wt=Math.clz32?Math.clz32:cc,ac=Math.log,uc=Math.LN2;function cc(e){return e>>>=0,e===0?32:31-(ac(e)/uc|0)|0}var el=64,tl=4194304;function dr(e){switch(e&-e){case 1:return 1;case 2:return 2;case 4:return 4;case 8:return 8;case 16:return 16;case 32:return 32;case 64:case 128:case 256:case 512:case 1024:case 2048:case 4096:case 8192:case 16384:case 32768:case 65536:case 131072:case 262144:case 524288:case 1048576:case 2097152:return e&4194240;case 4194304:case 8388608:case 16777216:case 33554432:case 67108864:return e&130023424;case 134217728:return 134217728;case 268435456:return 268435456;case 536870912:return 536870912;case 1073741824:return 1073741824;default:return e}}function nl(e,t){var n=e.pendingLanes;if(n===0)return 0;var r=0,l=e.suspendedLanes,o=e.pingedLanes,i=n&268435455;if(i!==0){var s=i&~l;s!==0?r=dr(s):(o&=i,o!==0&&(r=dr(o)))}else i=n&~l,i!==0?r=dr(i):o!==0&&(r=dr(o));if(r===0)return 0;if(t!==0&&t!==r&&(t&l)===0&&(l=r&-r,o=t&-t,l>=o||l===16&&(o&4194240)!==0))return t;if((r&4)!==0&&(r|=n&16),t=e.entangledLanes,t!==0)for(e=e.entanglements,t&=r;0<t;)n=31-wt(t),l=1<<n,r|=e[n],t&=~l;return r}function dc(e,t){switch(e){case 1:case 2:case 4:return t+250;case 8:case 16:case 32:case 64:case 128:case 256:case 512:case 1024:case 2048:case 4096:case 8192:case 16384:case 32768:case 65536:case 131072:case 262144:case 524288:case 1048576:case 2097152:return t+5e3;case 4194304:case 8388608:case 16777216:case 33554432:case 67108864:return-1;case 134217728:case 268435456:case 536870912:case 1073741824:return-1;default:return-1}}function fc(e,t){for(var n=e.suspendedLanes,r=e.pingedLanes,l=e.expirationTimes,o=e.pendingLanes;0<o;){var i=31-wt(o),s=1<<i,d=l[i];d===-1?((s&n)===0||(s&r)!==0)&&(l[i]=dc(s,t)):d<=t&&(e.expiredLanes|=s),o&=~s}}function mo(e){return e=e.pendingLanes&-1073741825,e!==0?e:e&1073741824?1073741824:0}function fs(){var e=el;return el<<=1,(el&4194240)===0&&(el=64),e}function ho(e){for(var t=[],n=0;31>n;n++)t.push(e);return t}function fr(e,t,n){e.pendingLanes|=t,t!==536870912&&(e.suspendedLanes=0,e.pingedLanes=0),e=e.eventTimes,t=31-wt(t),e[t]=n}function pc(e,t){var n=e.pendingLanes&~t;e.pendingLanes=t,e.suspendedLanes=0,e.pingedLanes=0,e.expiredLanes&=t,e.mutableReadLanes&=t,e.entangledLanes&=t,t=e.entanglements;var r=e.eventTimes;for(e=e.expirationTimes;0<n;){var l=31-wt(n),o=1<<l;t[l]=0,r[l]=-1,e[l]=-1,n&=~o}}function vo(e,t){var n=e.entangledLanes|=t;for(e=e.entanglements;n;){var r=31-wt(n),l=1<<r;l&t|e[r]&t&&(e[r]|=t),n&=~l}}var ae=0;function ps(e){return e&=-e,1<e?4<e?(e&268435455)!==0?16:536870912:4:1}var ms,go,hs,vs,gs,yo=!1,rl=[],Ht=null,Kt=null,Yt=null,pr=new Map,mr=new Map,Gt=[],mc="mousedown mouseup touchcancel touchend touchstart auxclick dblclick pointercancel pointerdown pointerup dragend dragstart drop compositionend compositionstart keydown keypress keyup input textInput copy cut paste click change contextmenu reset submit".split(" ");function ys(e,t){switch(e){case"focusin":case"focusout":Ht=null;break;case"dragenter":case"dragleave":Kt=null;break;case"mouseover":case"mouseout":Yt=null;break;case"pointerover":case"pointerout":pr.delete(t.pointerId);break;case"gotpointercapture":case"lostpointercapture":mr.delete(t.pointerId)}}function hr(e,t,n,r,l,o){return e===null||e.nativeEvent!==o?(e={blockedOn:t,domEventName:n,eventSystemFlags:r,nativeEvent:o,targetContainers:[l]},t!==null&&(t=Tr(t),t!==null&&go(t)),e):(e.eventSystemFlags|=r,t=e.targetContainers,l!==null&&t.indexOf(l)===-1&&t.push(l),e)}function hc(e,t,n,r,l){switch(t){case"focusin":return Ht=hr(Ht,e,t,n,r,l),!0;case"dragenter":return Kt=hr(Kt,e,t,n,r,l),!0;case"mouseover":return Yt=hr(Yt,e,t,n,r,l),!0;case"pointerover":var o=l.pointerId;return pr.set(o,hr(pr.get(o)||null,e,t,n,r,l)),!0;case"gotpointercapture":return o=l.pointerId,mr.set(o,hr(mr.get(o)||null,e,t,n,r,l)),!0}return!1}function _s(e){var t=wn(e.target);if(t!==null){var n=_n(t);if(n!==null){if(t=n.tag,t===13){if(t=ls(n),t!==null){e.blockedOn=t,gs(e.priority,function(){hs(n)});return}}else if(t===3&&n.stateNode.current.memoizedState.isDehydrated){e.blockedOn=n.tag===3?n.stateNode.containerInfo:null;return}}}e.blockedOn=null}function ll(e){if(e.blockedOn!==null)return!1;for(var t=e.targetContainers;0<t.length;){var n=wo(e.domEventName,e.eventSystemFlags,t[0],e.nativeEvent);if(n===null){n=e.nativeEvent;var r=new n.constructor(n.type,n);io=r,n.target.dispatchEvent(r),io=null}else return t=Tr(n),t!==null&&go(t),e.blockedOn=n,!1;t.shift()}return!0}function ws(e,t,n){ll(e)&&n.delete(t)}function vc(){yo=!1,Ht!==null&&ll(Ht)&&(Ht=null),Kt!==null&&ll(Kt)&&(Kt=null),Yt!==null&&ll(Yt)&&(Yt=null),pr.forEach(ws),mr.forEach(ws)}function vr(e,t){e.blockedOn===t&&(e.blockedOn=null,yo||(yo=!0,_.unstable_scheduleCallback(_.unstable_NormalPriority,vc)))}function gr(e){function t(l){return vr(l,e)}if(0<rl.length){vr(rl[0],e);for(var n=1;n<rl.length;n++){var r=rl[n];r.blockedOn===e&&(r.blockedOn=null)}}for(Ht!==null&&vr(Ht,e),Kt!==null&&vr(Kt,e),Yt!==null&&vr(Yt,e),pr.forEach(t),mr.forEach(t),n=0;n<Gt.length;n++)r=Gt[n],r.blockedOn===e&&(r.blockedOn=null);for(;0<Gt.length&&(n=Gt[0],n.blockedOn===null);)_s(n),n.blockedOn===null&&Gt.shift()}var Fn=Se.ReactCurrentBatchConfig,ol=!0;function gc(e,t,n,r){var l=ae,o=Fn.transition;Fn.transition=null;try{ae=1,_o(e,t,n,r)}finally{ae=l,Fn.transition=o}}function yc(e,t,n,r){var l=ae,o=Fn.transition;Fn.transition=null;try{ae=4,_o(e,t,n,r)}finally{ae=l,Fn.transition=o}}function _o(e,t,n,r){if(ol){var l=wo(e,t,n,r);if(l===null)Mo(e,t,r,il,n),ys(e,r);else if(hc(l,e,t,n,r))r.stopPropagation();else if(ys(e,r),t&4&&-1<mc.indexOf(e)){for(;l!==null;){var o=Tr(l);if(o!==null&&ms(o),o=wo(e,t,n,r),o===null&&Mo(e,t,r,il,n),o===l)break;l=o}l!==null&&r.stopPropagation()}else Mo(e,t,r,null,n)}}var il=null;function wo(e,t,n,r){if(il=null,e=so(r),e=wn(e),e!==null)if(t=_n(e),t===null)e=null;else if(n=t.tag,n===13){if(e=ls(t),e!==null)return e;e=null}else if(n===3){if(t.stateNode.current.memoizedState.isDehydrated)return t.tag===3?t.stateNode.containerInfo:null;e=null}else t!==e&&(e=null);return il=e,null}function ks(e){switch(e){case"cancel":case"click":case"close":case"contextmenu":case"copy":case"cut":case"auxclick":case"dblclick":case"dragend":case"dragstart":case"drop":case"focusin":case"focusout":case"input":case"invalid":case"keydown":case"keypress":case"keyup":case"mousedown":case"mouseup":case"paste":case"pause":case"play":case"pointercancel":case"pointerdown":case"pointerup":case"ratechange":case"reset":case"resize":case"seeked":case"submit":case"touchcancel":case"touchend":case"touchstart":case"volumechange":case"change":case"selectionchange":case"textInput":case"compositionstart":case"compositionend":case"compositionupdate":case"beforeblur":case"afterblur":case"beforeinput":case"blur":case"fullscreenchange":case"focus":case"hashchange":case"popstate":case"select":case"selectstart":return 1;case"drag":case"dragenter":case"dragexit":case"dragleave":case"dragover":case"mousemove":case"mouseout":case"mouseover":case"pointermove":case"pointerout":case"pointerover":case"scroll":case"toggle":case"touchmove":case"wheel":case"mouseenter":case"mouseleave":case"pointerenter":case"pointerleave":return 4;case"message":switch(oc()){case po:return 1;case cs:return 4;case Jr:case ic:return 16;case ds:return 536870912;default:return 16}default:return 16}}var Xt=null,ko=null,sl=null;function Ss(){if(sl)return sl;var e,t=ko,n=t.length,r,l="value"in Xt?Xt.value:Xt.textContent,o=l.length;for(e=0;e<n&&t[e]===l[e];e++);var i=n-e;for(r=1;r<=i&&t[n-r]===l[o-r];r++);return sl=l.slice(e,1<r?1-r:void 0)}function al(e){var t=e.keyCode;return"charCode"in e?(e=e.charCode,e===0&&t===13&&(e=13)):e=t,e===10&&(e=13),32<=e||e===13?e:0}function ul(){return!0}function $s(){return!1}function lt(e){function t(n,r,l,o,i){this._reactName=n,this._targetInst=l,this.type=r,this.nativeEvent=o,this.target=i,this.currentTarget=null;for(var s in e)e.hasOwnProperty(s)&&(n=e[s],this[s]=n?n(o):o[s]);return this.isDefaultPrevented=(o.defaultPrevented!=null?o.defaultPrevented:o.returnValue===!1)?ul:$s,this.isPropagationStopped=$s,this}return $(t.prototype,{preventDefault:function(){this.defaultPrevented=!0;var n=this.nativeEvent;n&&(n.preventDefault?n.preventDefault():typeof n.returnValue!="unknown"&&(n.returnValue=!1),this.isDefaultPrevented=ul)},stopPropagation:function(){var n=this.nativeEvent;n&&(n.stopPropagation?n.stopPropagation():typeof n.cancelBubble!="unknown"&&(n.cancelBubble=!0),this.isPropagationStopped=ul)},persist:function(){},isPersistent:ul}),t}var Mn={eventPhase:0,bubbles:0,cancelable:0,timeStamp:function(e){return e.timeStamp||Date.now()},defaultPrevented:0,isTrusted:0},So=lt(Mn),yr=$({},Mn,{view:0,detail:0}),_c=lt(yr),$o,No,_r,cl=$({},yr,{screenX:0,screenY:0,clientX:0,clientY:0,pageX:0,pageY:0,ctrlKey:0,shiftKey:0,altKey:0,metaKey:0,getModifierState:Eo,button:0,buttons:0,relatedTarget:function(e){return e.relatedTarget===void 0?e.fromElement===e.srcElement?e.toElement:e.fromElement:e.relatedTarget},movementX:function(e){return"movementX"in e?e.movementX:(e!==_r&&(_r&&e.type==="mousemove"?($o=e.screenX-_r.screenX,No=e.screenY-_r.screenY):No=$o=0,_r=e),$o)},movementY:function(e){return"movementY"in e?e.movementY:No}}),Ns=lt(cl),wc=$({},cl,{dataTransfer:0}),kc=lt(wc),Sc=$({},yr,{relatedTarget:0}),Co=lt(Sc),$c=$({},Mn,{animationName:0,elapsedTime:0,pseudoElement:0}),Nc=lt($c),Cc=$({},Mn,{clipboardData:function(e){return"clipboardData"in e?e.clipboardData:window.clipboardData}}),Ec=lt(Cc),xc=$({},Mn,{data:0}),Cs=lt(xc),Pc={Esc:"Escape",Spacebar:" ",Left:"ArrowLeft",Up:"ArrowUp",Right:"ArrowRight",Down:"ArrowDown",Del:"Delete",Win:"OS",Menu:"ContextMenu",Apps:"ContextMenu",Scroll:"ScrollLock",MozPrintableKey:"Unidentified"},Tc={8:"Backspace",9:"Tab",12:"Clear",13:"Enter",16:"Shift",17:"Control",18:"Alt",19:"Pause",20:"CapsLock",27:"Escape",32:" ",33:"PageUp",34:"PageDown",35:"End",36:"Home",37:"ArrowLeft",38:"ArrowUp",39:"ArrowRight",40:"ArrowDown",45:"Insert",46:"Delete",112:"F1",113:"F2",114:"F3",115:"F4",116:"F5",117:"F6",118:"F7",119:"F8",120:"F9",121:"F10",122:"F11",123:"F12",144:"NumLock",145:"ScrollLock",224:"Meta"},zc={Alt:"altKey",Control:"ctrlKey",Meta:"metaKey",Shift:"shiftKey"};function Rc(e){var t=this.nativeEvent;return t.getModifierState?t.getModifierState(e):(e=zc[e])?!!t[e]:!1}function Eo(){return Rc}var Lc=$({},yr,{key:function(e){if(e.key){var t=Pc[e.key]||e.key;if(t!=="Unidentified")return t}return e.type==="keypress"?(e=al(e),e===13?"Enter":String.fromCharCode(e)):e.type==="keydown"||e.type==="keyup"?Tc[e.keyCode]||"Unidentified":""},code:0,location:0,ctrlKey:0,shiftKey:0,altKey:0,metaKey:0,repeat:0,locale:0,getModifierState:Eo,charCode:function(e){return e.type==="keypress"?al(e):0},keyCode:function(e){return e.type==="keydown"||e.type==="keyup"?e.keyCode:0},which:function(e){return e.type==="keypress"?al(e):e.type==="keydown"||e.type==="keyup"?e.keyCode:0}}),Ic=lt(Lc),Dc=$({},cl,{pointerId:0,width:0,height:0,pressure:0,tangentialPressure:0,tiltX:0,tiltY:0,twist:0,pointerType:0,isPrimary:0}),Es=lt(Dc),Oc=$({},yr,{touches:0,targetTouches:0,changedTouches:0,altKey:0,metaKey:0,ctrlKey:0,shiftKey:0,getModifierState:Eo}),Fc=lt(Oc),Mc=$({},Mn,{propertyName:0,elapsedTime:0,pseudoElement:0}),jc=lt(Mc),Ac=$({},cl,{deltaX:function(e){return"deltaX"in e?e.deltaX:"wheelDeltaX"in e?-e.wheelDeltaX:0},deltaY:function(e){return"deltaY"in e?e.deltaY:"wheelDeltaY"in e?-e.wheelDeltaY:"wheelDelta"in e?-e.wheelDelta:0},deltaZ:0,deltaMode:0}),Uc=lt(Ac),Bc=[9,13,27,32],xo=G&&"CompositionEvent"in window,wr=null;G&&"documentMode"in document&&(wr=document.documentMode);var Wc=G&&"TextEvent"in window&&!wr,xs=G&&(!xo||wr&&8<wr&&11>=wr),Ps=" ",Ts=!1;function zs(e,t){switch(e){case"keyup":return Bc.indexOf(t.keyCode)!==-1;case"keydown":return t.keyCode!==229;case"keypress":case"mousedown":case"focusout":return!0;default:return!1}}function Rs(e){return e=e.detail,typeof e=="object"&&"data"in e?e.data:null}var jn=!1;function Qc(e,t){switch(e){case"compositionend":return Rs(t);case"keypress":return t.which!==32?null:(Ts=!0,Ps);case"textInput":return e=t.data,e===Ps&&Ts?null:e;default:return null}}function Vc(e,t){if(jn)return e==="compositionend"||!xo&&zs(e,t)?(e=Ss(),sl=ko=Xt=null,jn=!1,e):null;switch(e){case"paste":return null;case"keypress":if(!(t.ctrlKey||t.altKey||t.metaKey)||t.ctrlKey&&t.altKey){if(t.char&&1<t.char.length)return t.char;if(t.which)return String.fromCharCode(t.which)}return null;case"compositionend":return xs&&t.locale!=="ko"?null:t.data;default:return null}}var Hc={color:!0,date:!0,datetime:!0,"datetime-local":!0,email:!0,month:!0,number:!0,password:!0,range:!0,search:!0,tel:!0,text:!0,time:!0,url:!0,week:!0};function Ls(e){var t=e&&e.nodeName&&e.nodeName.toLowerCase();return t==="input"?!!Hc[e.type]:t==="textarea"}function Is(e,t,n,r){bi(r),t=hl(t,"onChange"),0<t.length&&(n=new So("onChange","change",null,n,r),e.push({event:n,listeners:t}))}var kr=null,Sr=null;function Kc(e){Zs(e,0)}function dl(e){var t=Qn(e);if(pn(t))return e}function Yc(e,t){if(e==="change")return t}var Ds=!1;if(G){var Po;if(G){var To="oninput"in document;if(!To){var Os=document.createElement("div");Os.setAttribute("oninput","return;"),To=typeof Os.oninput=="function"}Po=To}else Po=!1;Ds=Po&&(!document.documentMode||9<document.documentMode)}function Fs(){kr&&(kr.detachEvent("onpropertychange",Ms),Sr=kr=null)}function Ms(e){if(e.propertyName==="value"&&dl(Sr)){var t=[];Is(t,Sr,e,so(e)),rs(Kc,t)}}function Gc(e,t,n){e==="focusin"?(Fs(),kr=t,Sr=n,kr.attachEvent("onpropertychange",Ms)):e==="focusout"&&Fs()}function Xc(e){if(e==="selectionchange"||e==="keyup"||e==="keydown")return dl(Sr)}function qc(e,t){if(e==="click")return dl(t)}function Zc(e,t){if(e==="input"||e==="change")return dl(t)}function Jc(e,t){return e===t&&(e!==0||1/e===1/t)||e!==e&&t!==t}var kt=typeof Object.is=="function"?Object.is:Jc;function $r(e,t){if(kt(e,t))return!0;if(typeof e!="object"||e===null||typeof t!="object"||t===null)return!1;var n=Object.keys(e),r=Object.keys(t);if(n.length!==r.length)return!1;for(r=0;r<n.length;r++){var l=n[r];if(!C.call(t,l)||!kt(e[l],t[l]))return!1}return!0}function js(e){for(;e&&e.firstChild;)e=e.firstChild;return e}function As(e,t){var n=js(e);e=0;for(var r;n;){if(n.nodeType===3){if(r=e+n.textContent.length,e<=t&&r>=t)return{node:n,offset:t-e};e=r}e:{for(;n;){if(n.nextSibling){n=n.nextSibling;break e}n=n.parentNode}n=void 0}n=js(n)}}function Us(e,t){return e&&t?e===t?!0:e&&e.nodeType===3?!1:t&&t.nodeType===3?Us(e,t.parentNode):"contains"in e?e.contains(t):e.compareDocumentPosition?!!(e.compareDocumentPosition(t)&16):!1:!1}function Bs(){for(var e=window,t=mn();t instanceof e.HTMLIFrameElement;){try{var n=typeof t.contentWindow.location.href=="string"}catch{n=!1}if(n)e=t.contentWindow;else break;t=mn(e.document)}return t}function zo(e){var t=e&&e.nodeName&&e.nodeName.toLowerCase();return t&&(t==="input"&&(e.type==="text"||e.type==="search"||e.type==="tel"||e.type==="url"||e.type==="password")||t==="textarea"||e.contentEditable==="true")}function bc(e){var t=Bs(),n=e.focusedElem,r=e.selectionRange;if(t!==n&&n&&n.ownerDocument&&Us(n.ownerDocument.documentElement,n)){if(r!==null&&zo(n)){if(t=r.start,e=r.end,e===void 0&&(e=t),"selectionStart"in n)n.selectionStart=t,n.selectionEnd=Math.min(e,n.value.length);else if(e=(t=n.ownerDocument||document)&&t.defaultView||window,e.getSelection){e=e.getSelection();var l=n.textContent.length,o=Math.min(r.start,l);r=r.end===void 0?o:Math.min(r.end,l),!e.extend&&o>r&&(l=r,r=o,o=l),l=As(n,o);var i=As(n,r);l&&i&&(e.rangeCount!==1||e.anchorNode!==l.node||e.anchorOffset!==l.offset||e.focusNode!==i.node||e.focusOffset!==i.offset)&&(t=t.createRange(),t.setStart(l.node,l.offset),e.removeAllRanges(),o>r?(e.addRange(t),e.extend(i.node,i.offset)):(t.setEnd(i.node,i.offset),e.addRange(t)))}}for(t=[],e=n;e=e.parentNode;)e.nodeType===1&&t.push({element:e,left:e.scrollLeft,top:e.scrollTop});for(typeof n.focus=="function"&&n.focus(),n=0;n<t.length;n++)e=t[n],e.element.scrollLeft=e.left,e.element.scrollTop=e.top}}var ed=G&&"documentMode"in document&&11>=document.documentMode,An=null,Ro=null,Nr=null,Lo=!1;function Ws(e,t,n){var r=n.window===n?n.document:n.nodeType===9?n:n.ownerDocument;Lo||An==null||An!==mn(r)||(r=An,"selectionStart"in r&&zo(r)?r={start:r.selectionStart,end:r.selectionEnd}:(r=(r.ownerDocument&&r.ownerDocument.defaultView||window).getSelection(),r={anchorNode:r.anchorNode,anchorOffset:r.anchorOffset,focusNode:r.focusNode,focusOffset:r.focusOffset}),Nr&&$r(Nr,r)||(Nr=r,r=hl(Ro,"onSelect"),0<r.length&&(t=new So("onSelect","select",null,t,n),e.push({event:t,listeners:r}),t.target=An)))}function fl(e,t){var n={};return n[e.toLowerCase()]=t.toLowerCase(),n["Webkit"+e]="webkit"+t,n["Moz"+e]="moz"+t,n}var Un={animationend:fl("Animation","AnimationEnd"),animationiteration:fl("Animation","AnimationIteration"),animationstart:fl("Animation","AnimationStart"),transitionend:fl("Transition","TransitionEnd")},Io={},Qs={};G&&(Qs=document.createElement("div").style,"AnimationEvent"in window||(delete Un.animationend.animation,delete Un.animationiteration.animation,delete Un.animationstart.animation),"TransitionEvent"in window||delete Un.transitionend.transition);function pl(e){if(Io[e])return Io[e];if(!Un[e])return e;var t=Un[e],n;for(n in t)if(t.hasOwnProperty(n)&&n in Qs)return Io[e]=t[n];return e}var Vs=pl("animationend"),Hs=pl("animationiteration"),Ks=pl("animationstart"),Ys=pl("transitionend"),Gs=new Map,Xs="abort auxClick cancel canPlay canPlayThrough click close contextMenu copy cut drag dragEnd dragEnter dragExit dragLeave dragOver dragStart drop durationChange emptied encrypted ended error gotPointerCapture input invalid keyDown keyPress keyUp load loadedData loadedMetadata loadStart lostPointerCapture mouseDown mouseMove mouseOut mouseOver mouseUp paste pause play playing pointerCancel pointerDown pointerMove pointerOut pointerOver pointerUp progress rateChange reset resize seeked seeking stalled submit suspend timeUpdate touchCancel touchEnd touchStart volumeChange scroll toggle touchMove waiting wheel".split(" ");function qt(e,t){Gs.set(e,t),x(t,[e])}for(var Do=0;Do<Xs.length;Do++){var Oo=Xs[Do],td=Oo.toLowerCase(),nd=Oo[0].toUpperCase()+Oo.slice(1);qt(td,"on"+nd)}qt(Vs,"onAnimationEnd"),qt(Hs,"onAnimationIteration"),qt(Ks,"onAnimationStart"),qt("dblclick","onDoubleClick"),qt("focusin","onFocus"),qt("focusout","onBlur"),qt(Ys,"onTransitionEnd"),M("onMouseEnter",["mouseout","mouseover"]),M("onMouseLeave",["mouseout","mouseover"]),M("onPointerEnter",["pointerout","pointerover"]),M("onPointerLeave",["pointerout","pointerover"]),x("onChange","change click focusin focusout input keydown keyup selectionchange".split(" ")),x("onSelect","focusout contextmenu dragend focusin keydown keyup mousedown mouseup selectionchange".split(" ")),x("onBeforeInput",["compositionend","keypress","textInput","paste"]),x("onCompositionEnd","compositionend focusout keydown keypress keyup mousedown".split(" ")),x("onCompositionStart","compositionstart focusout keydown keypress keyup mousedown".split(" ")),x("onCompositionUpdate","compositionupdate focusout keydown keypress keyup mousedown".split(" "));var Cr="abort canplay canplaythrough durationchange emptied encrypted ended error loadeddata loadedmetadata loadstart pause play playing progress ratechange resize seeked seeking stalled suspend timeupdate volumechange waiting".split(" "),rd=new Set("cancel close invalid load scroll toggle".split(" ").concat(Cr));function qs(e,t,n){var r=e.type||"unknown-event";e.currentTarget=n,tc(r,t,void 0,e),e.currentTarget=null}function Zs(e,t){t=(t&4)!==0;for(var n=0;n<e.length;n++){var r=e[n],l=r.event;r=r.listeners;e:{var o=void 0;if(t)for(var i=r.length-1;0<=i;i--){var s=r[i],d=s.instance,g=s.currentTarget;if(s=s.listener,d!==o&&l.isPropagationStopped())break e;qs(l,s,g),o=d}else for(i=0;i<r.length;i++){if(s=r[i],d=s.instance,g=s.currentTarget,s=s.listener,d!==o&&l.isPropagationStopped())break e;qs(l,s,g),o=d}}}if(Zr)throw e=fo,Zr=!1,fo=null,e}function he(e,t){var n=t[Qo];n===void 0&&(n=t[Qo]=new Set);var r=e+"__bubble";n.has(r)||(Js(t,e,2,!1),n.add(r))}function Fo(e,t,n){var r=0;t&&(r|=4),Js(n,e,r,t)}var ml="_reactListening"+Math.random().toString(36).slice(2);function Er(e){if(!e[ml]){e[ml]=!0,T.forEach(function(n){n!=="selectionchange"&&(rd.has(n)||Fo(n,!1,e),Fo(n,!0,e))});var t=e.nodeType===9?e:e.ownerDocument;t===null||t[ml]||(t[ml]=!0,Fo("selectionchange",!1,t))}}function Js(e,t,n,r){switch(ks(t)){case 1:var l=gc;break;case 4:l=yc;break;default:l=_o}n=l.bind(null,t,n,e),l=void 0,!co||t!=="touchstart"&&t!=="touchmove"&&t!=="wheel"||(l=!0),r?l!==void 0?e.addEventListener(t,n,{capture:!0,passive:l}):e.addEventListener(t,n,!0):l!==void 0?e.addEventListener(t,n,{passive:l}):e.addEventListener(t,n,!1)}function Mo(e,t,n,r,l){var o=r;if((t&1)===0&&(t&2)===0&&r!==null)e:for(;;){if(r===null)return;var i=r.tag;if(i===3||i===4){var s=r.stateNode.containerInfo;if(s===l||s.nodeType===8&&s.parentNode===l)break;if(i===4)for(i=r.return;i!==null;){var d=i.tag;if((d===3||d===4)&&(d=i.stateNode.containerInfo,d===l||d.nodeType===8&&d.parentNode===l))return;i=i.return}for(;s!==null;){if(i=wn(s),i===null)return;if(d=i.tag,d===5||d===6){r=o=i;continue e}s=s.parentNode}}r=r.return}rs(function(){var g=o,k=so(n),S=[];e:{var w=Gs.get(e);if(w!==void 0){var R=So,O=e;switch(e){case"keypress":if(al(n)===0)break e;case"keydown":case"keyup":R=Ic;break;case"focusin":O="focus",R=Co;break;case"focusout":O="blur",R=Co;break;case"beforeblur":case"afterblur":R=Co;break;case"click":if(n.button===2)break e;case"auxclick":case"dblclick":case"mousedown":case"mousemove":case"mouseup":case"mouseout":case"mouseover":case"contextmenu":R=Ns;break;case"drag":case"dragend":case"dragenter":case"dragexit":case"dragleave":case"dragover":case"dragstart":case"drop":R=kc;break;case"touchcancel":case"touchend":case"touchmove":case"touchstart":R=Fc;break;case Vs:case Hs:case Ks:R=Nc;break;case Ys:R=jc;break;case"scroll":R=_c;break;case"wheel":R=Uc;break;case"copy":case"cut":case"paste":R=Ec;break;case"gotpointercapture":case"lostpointercapture":case"pointercancel":case"pointerdown":case"pointermove":case"pointerout":case"pointerover":case"pointerup":R=Es}var F=(t&4)!==0,Ne=!F&&e==="scroll",h=F?w!==null?w+"Capture":null:w;F=[];for(var p=g,v;p!==null;){v=p;var N=v.stateNode;if(v.tag===5&&N!==null&&(v=N,h!==null&&(N=ar(p,h),N!=null&&F.push(xr(p,N,v)))),Ne)break;p=p.return}0<F.length&&(w=new R(w,O,null,n,k),S.push({event:w,listeners:F}))}}if((t&7)===0){e:{if(w=e==="mouseover"||e==="pointerover",R=e==="mouseout"||e==="pointerout",w&&n!==io&&(O=n.relatedTarget||n.fromElement)&&(wn(O)||O[Dt]))break e;if((R||w)&&(w=k.window===k?k:(w=k.ownerDocument)?w.defaultView||w.parentWindow:window,R?(O=n.relatedTarget||n.toElement,R=g,O=O?wn(O):null,O!==null&&(Ne=_n(O),O!==Ne||O.tag!==5&&O.tag!==6)&&(O=null)):(R=null,O=g),R!==O)){if(F=Ns,N="onMouseLeave",h="onMouseEnter",p="mouse",(e==="pointerout"||e==="pointerover")&&(F=Es,N="onPointerLeave",h="onPointerEnter",p="pointer"),Ne=R==null?w:Qn(R),v=O==null?w:Qn(O),w=new F(N,p+"leave",R,n,k),w.target=Ne,w.relatedTarget=v,N=null,wn(k)===g&&(F=new F(h,p+"enter",O,n,k),F.target=v,F.relatedTarget=Ne,N=F),Ne=N,R&&O)t:{for(F=R,h=O,p=0,v=F;v;v=Bn(v))p++;for(v=0,N=h;N;N=Bn(N))v++;for(;0<p-v;)F=Bn(F),p--;for(;0<v-p;)h=Bn(h),v--;for(;p--;){if(F===h||h!==null&&F===h.alternate)break t;F=Bn(F),h=Bn(h)}F=null}else F=null;R!==null&&bs(S,w,R,F,!1),O!==null&&Ne!==null&&bs(S,Ne,O,F,!0)}}e:{if(w=g?Qn(g):window,R=w.nodeName&&w.nodeName.toLowerCase(),R==="select"||R==="input"&&w.type==="file")var A=Yc;else if(Ls(w))if(Ds)A=Zc;else{A=Xc;var B=Gc}else(R=w.nodeName)&&R.toLowerCase()==="input"&&(w.type==="checkbox"||w.type==="radio")&&(A=qc);if(A&&(A=A(e,g))){Is(S,A,n,k);break e}B&&B(e,w,g),e==="focusout"&&(B=w._wrapperState)&&B.controlled&&w.type==="number"&&ir(w,"number",w.value)}switch(B=g?Qn(g):window,e){case"focusin":(Ls(B)||B.contentEditable==="true")&&(An=B,Ro=g,Nr=null);break;case"focusout":Nr=Ro=An=null;break;case"mousedown":Lo=!0;break;case"contextmenu":case"mouseup":case"dragend":Lo=!1,Ws(S,n,k);break;case"selectionchange":if(ed)break;case"keydown":case"keyup":Ws(S,n,k)}var W;if(xo)e:{switch(e){case"compositionstart":var H="onCompositionStart";break e;case"compositionend":H="onCompositionEnd";break e;case"compositionupdate":H="onCompositionUpdate";break e}H=void 0}else jn?zs(e,n)&&(H="onCompositionEnd"):e==="keydown"&&n.keyCode===229&&(H="onCompositionStart");H&&(xs&&n.locale!=="ko"&&(jn||H!=="onCompositionStart"?H==="onCompositionEnd"&&jn&&(W=Ss()):(Xt=k,ko="value"in Xt?Xt.value:Xt.textContent,jn=!0)),B=hl(g,H),0<B.length&&(H=new Cs(H,e,null,n,k),S.push({event:H,listeners:B}),W?H.data=W:(W=Rs(n),W!==null&&(H.data=W)))),(W=Wc?Qc(e,n):Vc(e,n))&&(g=hl(g,"onBeforeInput"),0<g.length&&(k=new Cs("onBeforeInput","beforeinput",null,n,k),S.push({event:k,listeners:g}),k.data=W))}Zs(S,t)})}function xr(e,t,n){return{instance:e,listener:t,currentTarget:n}}function hl(e,t){for(var n=t+"Capture",r=[];e!==null;){var l=e,o=l.stateNode;l.tag===5&&o!==null&&(l=o,o=ar(e,n),o!=null&&r.unshift(xr(e,o,l)),o=ar(e,t),o!=null&&r.push(xr(e,o,l))),e=e.return}return r}function Bn(e){if(e===null)return null;do e=e.return;while(e&&e.tag!==5);return e||null}function bs(e,t,n,r,l){for(var o=t._reactName,i=[];n!==null&&n!==r;){var s=n,d=s.alternate,g=s.stateNode;if(d!==null&&d===r)break;s.tag===5&&g!==null&&(s=g,l?(d=ar(n,o),d!=null&&i.unshift(xr(n,d,s))):l||(d=ar(n,o),d!=null&&i.push(xr(n,d,s)))),n=n.return}i.length!==0&&e.push({event:t,listeners:i})}var ld=/\r\n?/g,od=/\u0000|\uFFFD/g;function ea(e){return(typeof e=="string"?e:""+e).replace(ld,`
`).replace(od,"")}function vl(e,t,n){if(t=ea(t),ea(e)!==t&&n)throw Error(m(425))}function gl(){}var jo=null,Ao=null;function Uo(e,t){return e==="textarea"||e==="noscript"||typeof t.children=="string"||typeof t.children=="number"||typeof t.dangerouslySetInnerHTML=="object"&&t.dangerouslySetInnerHTML!==null&&t.dangerouslySetInnerHTML.__html!=null}var Bo=typeof setTimeout=="function"?setTimeout:void 0,id=typeof clearTimeout=="function"?clearTimeout:void 0,ta=typeof Promise=="function"?Promise:void 0,sd=typeof queueMicrotask=="function"?queueMicrotask:typeof ta<"u"?function(e){return ta.resolve(null).then(e).catch(ad)}:Bo;function ad(e){setTimeout(function(){throw e})}function Wo(e,t){var n=t,r=0;do{var l=n.nextSibling;if(e.removeChild(n),l&&l.nodeType===8)if(n=l.data,n==="/$"){if(r===0){e.removeChild(l),gr(t);return}r--}else n!=="$"&&n!=="$?"&&n!=="$!"||r++;n=l}while(n);gr(t)}function Zt(e){for(;e!=null;e=e.nextSibling){var t=e.nodeType;if(t===1||t===3)break;if(t===8){if(t=e.data,t==="$"||t==="$!"||t==="$?")break;if(t==="/$")return null}}return e}function na(e){e=e.previousSibling;for(var t=0;e;){if(e.nodeType===8){var n=e.data;if(n==="$"||n==="$!"||n==="$?"){if(t===0)return e;t--}else n==="/$"&&t++}e=e.previousSibling}return null}var Wn=Math.random().toString(36).slice(2),Tt="__reactFiber$"+Wn,Pr="__reactProps$"+Wn,Dt="__reactContainer$"+Wn,Qo="__reactEvents$"+Wn,ud="__reactListeners$"+Wn,cd="__reactHandles$"+Wn;function wn(e){var t=e[Tt];if(t)return t;for(var n=e.parentNode;n;){if(t=n[Dt]||n[Tt]){if(n=t.alternate,t.child!==null||n!==null&&n.child!==null)for(e=na(e);e!==null;){if(n=e[Tt])return n;e=na(e)}return t}e=n,n=e.parentNode}return null}function Tr(e){return e=e[Tt]||e[Dt],!e||e.tag!==5&&e.tag!==6&&e.tag!==13&&e.tag!==3?null:e}function Qn(e){if(e.tag===5||e.tag===6)return e.stateNode;throw Error(m(33))}function yl(e){return e[Pr]||null}var Vo=[],Vn=-1;function Jt(e){return{current:e}}function ve(e){0>Vn||(e.current=Vo[Vn],Vo[Vn]=null,Vn--)}function pe(e,t){Vn++,Vo[Vn]=e.current,e.current=t}var bt={},Qe=Jt(bt),Xe=Jt(!1),kn=bt;function Hn(e,t){var n=e.type.contextTypes;if(!n)return bt;var r=e.stateNode;if(r&&r.__reactInternalMemoizedUnmaskedChildContext===t)return r.__reactInternalMemoizedMaskedChildContext;var l={},o;for(o in n)l[o]=t[o];return r&&(e=e.stateNode,e.__reactInternalMemoizedUnmaskedChildContext=t,e.__reactInternalMemoizedMaskedChildContext=l),l}function qe(e){return e=e.childContextTypes,e!=null}function _l(){ve(Xe),ve(Qe)}function ra(e,t,n){if(Qe.current!==bt)throw Error(m(168));pe(Qe,t),pe(Xe,n)}function la(e,t,n){var r=e.stateNode;if(t=t.childContextTypes,typeof r.getChildContext!="function")return n;r=r.getChildContext();for(var l in r)if(!(l in t))throw Error(m(108,b(e)||"Unknown",l));return $({},n,r)}function wl(e){return e=(e=e.stateNode)&&e.__reactInternalMemoizedMergedChildContext||bt,kn=Qe.current,pe(Qe,e),pe(Xe,Xe.current),!0}function oa(e,t,n){var r=e.stateNode;if(!r)throw Error(m(169));n?(e=la(e,t,kn),r.__reactInternalMemoizedMergedChildContext=e,ve(Xe),ve(Qe),pe(Qe,e)):ve(Xe),pe(Xe,n)}var Ot=null,kl=!1,Ho=!1;function ia(e){Ot===null?Ot=[e]:Ot.push(e)}function dd(e){kl=!0,ia(e)}function en(){if(!Ho&&Ot!==null){Ho=!0;var e=0,t=ae;try{var n=Ot;for(ae=1;e<n.length;e++){var r=n[e];do r=r(!0);while(r!==null)}Ot=null,kl=!1}catch(l){throw Ot!==null&&(Ot=Ot.slice(e+1)),as(po,en),l}finally{ae=t,Ho=!1}}return null}var Kn=[],Yn=0,Sl=null,$l=0,ct=[],dt=0,Sn=null,Ft=1,Mt="";function $n(e,t){Kn[Yn++]=$l,Kn[Yn++]=Sl,Sl=e,$l=t}function sa(e,t,n){ct[dt++]=Ft,ct[dt++]=Mt,ct[dt++]=Sn,Sn=e;var r=Ft;e=Mt;var l=32-wt(r)-1;r&=~(1<<l),n+=1;var o=32-wt(t)+l;if(30<o){var i=l-l%5;o=(r&(1<<i)-1).toString(32),r>>=i,l-=i,Ft=1<<32-wt(t)+l|n<<l|r,Mt=o+e}else Ft=1<<o|n<<l|r,Mt=e}function Ko(e){e.return!==null&&($n(e,1),sa(e,1,0))}function Yo(e){for(;e===Sl;)Sl=Kn[--Yn],Kn[Yn]=null,$l=Kn[--Yn],Kn[Yn]=null;for(;e===Sn;)Sn=ct[--dt],ct[dt]=null,Mt=ct[--dt],ct[dt]=null,Ft=ct[--dt],ct[dt]=null}var ot=null,it=null,ge=!1,St=null;function aa(e,t){var n=ht(5,null,null,0);n.elementType="DELETED",n.stateNode=t,n.return=e,t=e.deletions,t===null?(e.deletions=[n],e.flags|=16):t.push(n)}function ua(e,t){switch(e.tag){case 5:var n=e.type;return t=t.nodeType!==1||n.toLowerCase()!==t.nodeName.toLowerCase()?null:t,t!==null?(e.stateNode=t,ot=e,it=Zt(t.firstChild),!0):!1;case 6:return t=e.pendingProps===""||t.nodeType!==3?null:t,t!==null?(e.stateNode=t,ot=e,it=null,!0):!1;case 13:return t=t.nodeType!==8?null:t,t!==null?(n=Sn!==null?{id:Ft,overflow:Mt}:null,e.memoizedState={dehydrated:t,treeContext:n,retryLane:1073741824},n=ht(18,null,null,0),n.stateNode=t,n.return=e,e.child=n,ot=e,it=null,!0):!1;default:return!1}}function Go(e){return(e.mode&1)!==0&&(e.flags&128)===0}function Xo(e){if(ge){var t=it;if(t){var n=t;if(!ua(e,t)){if(Go(e))throw Error(m(418));t=Zt(n.nextSibling);var r=ot;t&&ua(e,t)?aa(r,n):(e.flags=e.flags&-4097|2,ge=!1,ot=e)}}else{if(Go(e))throw Error(m(418));e.flags=e.flags&-4097|2,ge=!1,ot=e}}}function ca(e){for(e=e.return;e!==null&&e.tag!==5&&e.tag!==3&&e.tag!==13;)e=e.return;ot=e}function Nl(e){if(e!==ot)return!1;if(!ge)return ca(e),ge=!0,!1;var t;if((t=e.tag!==3)&&!(t=e.tag!==5)&&(t=e.type,t=t!=="head"&&t!=="body"&&!Uo(e.type,e.memoizedProps)),t&&(t=it)){if(Go(e))throw da(),Error(m(418));for(;t;)aa(e,t),t=Zt(t.nextSibling)}if(ca(e),e.tag===13){if(e=e.memoizedState,e=e!==null?e.dehydrated:null,!e)throw Error(m(317));e:{for(e=e.nextSibling,t=0;e;){if(e.nodeType===8){var n=e.data;if(n==="/$"){if(t===0){it=Zt(e.nextSibling);break e}t--}else n!=="$"&&n!=="$!"&&n!=="$?"||t++}e=e.nextSibling}it=null}}else it=ot?Zt(e.stateNode.nextSibling):null;return!0}function da(){for(var e=it;e;)e=Zt(e.nextSibling)}function Gn(){it=ot=null,ge=!1}function qo(e){St===null?St=[e]:St.push(e)}var fd=Se.ReactCurrentBatchConfig;function zr(e,t,n){if(e=n.ref,e!==null&&typeof e!="function"&&typeof e!="object"){if(n._owner){if(n=n._owner,n){if(n.tag!==1)throw Error(m(309));var r=n.stateNode}if(!r)throw Error(m(147,e));var l=r,o=""+e;return t!==null&&t.ref!==null&&typeof t.ref=="function"&&t.ref._stringRef===o?t.ref:(t=function(i){var s=l.refs;i===null?delete s[o]:s[o]=i},t._stringRef=o,t)}if(typeof e!="string")throw Error(m(284));if(!n._owner)throw Error(m(290,e))}return e}function Cl(e,t){throw e=Object.prototype.toString.call(t),Error(m(31,e==="[object Object]"?"object with keys {"+Object.keys(t).join(", ")+"}":e))}function fa(e){var t=e._init;return t(e._payload)}function pa(e){function t(h,p){if(e){var v=h.deletions;v===null?(h.deletions=[p],h.flags|=16):v.push(p)}}function n(h,p){if(!e)return null;for(;p!==null;)t(h,p),p=p.sibling;return null}function r(h,p){for(h=new Map;p!==null;)p.key!==null?h.set(p.key,p):h.set(p.index,p),p=p.sibling;return h}function l(h,p){return h=un(h,p),h.index=0,h.sibling=null,h}function o(h,p,v){return h.index=v,e?(v=h.alternate,v!==null?(v=v.index,v<p?(h.flags|=2,p):v):(h.flags|=2,p)):(h.flags|=1048576,p)}function i(h){return e&&h.alternate===null&&(h.flags|=2),h}function s(h,p,v,N){return p===null||p.tag!==6?(p=Bi(v,h.mode,N),p.return=h,p):(p=l(p,v),p.return=h,p)}function d(h,p,v,N){var A=v.type;return A===Me?k(h,p,v.props.children,N,v.key):p!==null&&(p.elementType===A||typeof A=="object"&&A!==null&&A.$$typeof===j&&fa(A)===p.type)?(N=l(p,v.props),N.ref=zr(h,p,v),N.return=h,N):(N=Xl(v.type,v.key,v.props,null,h.mode,N),N.ref=zr(h,p,v),N.return=h,N)}function g(h,p,v,N){return p===null||p.tag!==4||p.stateNode.containerInfo!==v.containerInfo||p.stateNode.implementation!==v.implementation?(p=Wi(v,h.mode,N),p.return=h,p):(p=l(p,v.children||[]),p.return=h,p)}function k(h,p,v,N,A){return p===null||p.tag!==7?(p=Rn(v,h.mode,N,A),p.return=h,p):(p=l(p,v),p.return=h,p)}function S(h,p,v){if(typeof p=="string"&&p!==""||typeof p=="number")return p=Bi(""+p,h.mode,v),p.return=h,p;if(typeof p=="object"&&p!==null){switch(p.$$typeof){case rt:return v=Xl(p.type,p.key,p.props,null,h.mode,v),v.ref=zr(h,null,p),v.return=h,v;case Fe:return p=Wi(p,h.mode,v),p.return=h,p;case j:var N=p._init;return S(h,N(p._payload),v)}if(vn(p)||U(p))return p=Rn(p,h.mode,v,null),p.return=h,p;Cl(h,p)}return null}function w(h,p,v,N){var A=p!==null?p.key:null;if(typeof v=="string"&&v!==""||typeof v=="number")return A!==null?null:s(h,p,""+v,N);if(typeof v=="object"&&v!==null){switch(v.$$typeof){case rt:return v.key===A?d(h,p,v,N):null;case Fe:return v.key===A?g(h,p,v,N):null;case j:return A=v._init,w(h,p,A(v._payload),N)}if(vn(v)||U(v))return A!==null?null:k(h,p,v,N,null);Cl(h,v)}return null}function R(h,p,v,N,A){if(typeof N=="string"&&N!==""||typeof N=="number")return h=h.get(v)||null,s(p,h,""+N,A);if(typeof N=="object"&&N!==null){switch(N.$$typeof){case rt:return h=h.get(N.key===null?v:N.key)||null,d(p,h,N,A);case Fe:return h=h.get(N.key===null?v:N.key)||null,g(p,h,N,A);case j:var B=N._init;return R(h,p,v,B(N._payload),A)}if(vn(N)||U(N))return h=h.get(v)||null,k(p,h,N,A,null);Cl(p,N)}return null}function O(h,p,v,N){for(var A=null,B=null,W=p,H=p=0,Ie=null;W!==null&&H<v.length;H++){W.index>H?(Ie=W,W=null):Ie=W.sibling;var ne=w(h,W,v[H],N);if(ne===null){W===null&&(W=Ie);break}e&&W&&ne.alternate===null&&t(h,W),p=o(ne,p,H),B===null?A=ne:B.sibling=ne,B=ne,W=Ie}if(H===v.length)return n(h,W),ge&&$n(h,H),A;if(W===null){for(;H<v.length;H++)W=S(h,v[H],N),W!==null&&(p=o(W,p,H),B===null?A=W:B.sibling=W,B=W);return ge&&$n(h,H),A}for(W=r(h,W);H<v.length;H++)Ie=R(W,h,H,v[H],N),Ie!==null&&(e&&Ie.alternate!==null&&W.delete(Ie.key===null?H:Ie.key),p=o(Ie,p,H),B===null?A=Ie:B.sibling=Ie,B=Ie);return e&&W.forEach(function(cn){return t(h,cn)}),ge&&$n(h,H),A}function F(h,p,v,N){var A=U(v);if(typeof A!="function")throw Error(m(150));if(v=A.call(v),v==null)throw Error(m(151));for(var B=A=null,W=p,H=p=0,Ie=null,ne=v.next();W!==null&&!ne.done;H++,ne=v.next()){W.index>H?(Ie=W,W=null):Ie=W.sibling;var cn=w(h,W,ne.value,N);if(cn===null){W===null&&(W=Ie);break}e&&W&&cn.alternate===null&&t(h,W),p=o(cn,p,H),B===null?A=cn:B.sibling=cn,B=cn,W=Ie}if(ne.done)return n(h,W),ge&&$n(h,H),A;if(W===null){for(;!ne.done;H++,ne=v.next())ne=S(h,ne.value,N),ne!==null&&(p=o(ne,p,H),B===null?A=ne:B.sibling=ne,B=ne);return ge&&$n(h,H),A}for(W=r(h,W);!ne.done;H++,ne=v.next())ne=R(W,h,H,ne.value,N),ne!==null&&(e&&ne.alternate!==null&&W.delete(ne.key===null?H:ne.key),p=o(ne,p,H),B===null?A=ne:B.sibling=ne,B=ne);return e&&W.forEach(function(Hd){return t(h,Hd)}),ge&&$n(h,H),A}function Ne(h,p,v,N){if(typeof v=="object"&&v!==null&&v.type===Me&&v.key===null&&(v=v.props.children),typeof v=="object"&&v!==null){switch(v.$$typeof){case rt:e:{for(var A=v.key,B=p;B!==null;){if(B.key===A){if(A=v.type,A===Me){if(B.tag===7){n(h,B.sibling),p=l(B,v.props.children),p.return=h,h=p;break e}}else if(B.elementType===A||typeof A=="object"&&A!==null&&A.$$typeof===j&&fa(A)===B.type){n(h,B.sibling),p=l(B,v.props),p.ref=zr(h,B,v),p.return=h,h=p;break e}n(h,B);break}else t(h,B);B=B.sibling}v.type===Me?(p=Rn(v.props.children,h.mode,N,v.key),p.return=h,h=p):(N=Xl(v.type,v.key,v.props,null,h.mode,N),N.ref=zr(h,p,v),N.return=h,h=N)}return i(h);case Fe:e:{for(B=v.key;p!==null;){if(p.key===B)if(p.tag===4&&p.stateNode.containerInfo===v.containerInfo&&p.stateNode.implementation===v.implementation){n(h,p.sibling),p=l(p,v.children||[]),p.return=h,h=p;break e}else{n(h,p);break}else t(h,p);p=p.sibling}p=Wi(v,h.mode,N),p.return=h,h=p}return i(h);case j:return B=v._init,Ne(h,p,B(v._payload),N)}if(vn(v))return O(h,p,v,N);if(U(v))return F(h,p,v,N);Cl(h,v)}return typeof v=="string"&&v!==""||typeof v=="number"?(v=""+v,p!==null&&p.tag===6?(n(h,p.sibling),p=l(p,v),p.return=h,h=p):(n(h,p),p=Bi(v,h.mode,N),p.return=h,h=p),i(h)):n(h,p)}return Ne}var Xn=pa(!0),ma=pa(!1),El=Jt(null),xl=null,qn=null,Zo=null;function Jo(){Zo=qn=xl=null}function bo(e){var t=El.current;ve(El),e._currentValue=t}function ei(e,t,n){for(;e!==null;){var r=e.alternate;if((e.childLanes&t)!==t?(e.childLanes|=t,r!==null&&(r.childLanes|=t)):r!==null&&(r.childLanes&t)!==t&&(r.childLanes|=t),e===n)break;e=e.return}}function Zn(e,t){xl=e,Zo=qn=null,e=e.dependencies,e!==null&&e.firstContext!==null&&((e.lanes&t)!==0&&(Ze=!0),e.firstContext=null)}function ft(e){var t=e._currentValue;if(Zo!==e)if(e={context:e,memoizedValue:t,next:null},qn===null){if(xl===null)throw Error(m(308));qn=e,xl.dependencies={lanes:0,firstContext:e}}else qn=qn.next=e;return t}var Nn=null;function ti(e){Nn===null?Nn=[e]:Nn.push(e)}function ha(e,t,n,r){var l=t.interleaved;return l===null?(n.next=n,ti(t)):(n.next=l.next,l.next=n),t.interleaved=n,jt(e,r)}function jt(e,t){e.lanes|=t;var n=e.alternate;for(n!==null&&(n.lanes|=t),n=e,e=e.return;e!==null;)e.childLanes|=t,n=e.alternate,n!==null&&(n.childLanes|=t),n=e,e=e.return;return n.tag===3?n.stateNode:null}var tn=!1;function ni(e){e.updateQueue={baseState:e.memoizedState,firstBaseUpdate:null,lastBaseUpdate:null,shared:{pending:null,interleaved:null,lanes:0},effects:null}}function va(e,t){e=e.updateQueue,t.updateQueue===e&&(t.updateQueue={baseState:e.baseState,firstBaseUpdate:e.firstBaseUpdate,lastBaseUpdate:e.lastBaseUpdate,shared:e.shared,effects:e.effects})}function At(e,t){return{eventTime:e,lane:t,tag:0,payload:null,callback:null,next:null}}function nn(e,t,n){var r=e.updateQueue;if(r===null)return null;if(r=r.shared,(te&2)!==0){var l=r.pending;return l===null?t.next=t:(t.next=l.next,l.next=t),r.pending=t,jt(e,n)}return l=r.interleaved,l===null?(t.next=t,ti(r)):(t.next=l.next,l.next=t),r.interleaved=t,jt(e,n)}function Pl(e,t,n){if(t=t.updateQueue,t!==null&&(t=t.shared,(n&4194240)!==0)){var r=t.lanes;r&=e.pendingLanes,n|=r,t.lanes=n,vo(e,n)}}function ga(e,t){var n=e.updateQueue,r=e.alternate;if(r!==null&&(r=r.updateQueue,n===r)){var l=null,o=null;if(n=n.firstBaseUpdate,n!==null){do{var i={eventTime:n.eventTime,lane:n.lane,tag:n.tag,payload:n.payload,callback:n.callback,next:null};o===null?l=o=i:o=o.next=i,n=n.next}while(n!==null);o===null?l=o=t:o=o.next=t}else l=o=t;n={baseState:r.baseState,firstBaseUpdate:l,lastBaseUpdate:o,shared:r.shared,effects:r.effects},e.updateQueue=n;return}e=n.lastBaseUpdate,e===null?n.firstBaseUpdate=t:e.next=t,n.lastBaseUpdate=t}function Tl(e,t,n,r){var l=e.updateQueue;tn=!1;var o=l.firstBaseUpdate,i=l.lastBaseUpdate,s=l.shared.pending;if(s!==null){l.shared.pending=null;var d=s,g=d.next;d.next=null,i===null?o=g:i.next=g,i=d;var k=e.alternate;k!==null&&(k=k.updateQueue,s=k.lastBaseUpdate,s!==i&&(s===null?k.firstBaseUpdate=g:s.next=g,k.lastBaseUpdate=d))}if(o!==null){var S=l.baseState;i=0,k=g=d=null,s=o;do{var w=s.lane,R=s.eventTime;if((r&w)===w){k!==null&&(k=k.next={eventTime:R,lane:0,tag:s.tag,payload:s.payload,callback:s.callback,next:null});e:{var O=e,F=s;switch(w=t,R=n,F.tag){case 1:if(O=F.payload,typeof O=="function"){S=O.call(R,S,w);break e}S=O;break e;case 3:O.flags=O.flags&-65537|128;case 0:if(O=F.payload,w=typeof O=="function"?O.call(R,S,w):O,w==null)break e;S=$({},S,w);break e;case 2:tn=!0}}s.callback!==null&&s.lane!==0&&(e.flags|=64,w=l.effects,w===null?l.effects=[s]:w.push(s))}else R={eventTime:R,lane:w,tag:s.tag,payload:s.payload,callback:s.callback,next:null},k===null?(g=k=R,d=S):k=k.next=R,i|=w;if(s=s.next,s===null){if(s=l.shared.pending,s===null)break;w=s,s=w.next,w.next=null,l.lastBaseUpdate=w,l.shared.pending=null}}while(!0);if(k===null&&(d=S),l.baseState=d,l.firstBaseUpdate=g,l.lastBaseUpdate=k,t=l.shared.interleaved,t!==null){l=t;do i|=l.lane,l=l.next;while(l!==t)}else o===null&&(l.shared.lanes=0);xn|=i,e.lanes=i,e.memoizedState=S}}function ya(e,t,n){if(e=t.effects,t.effects=null,e!==null)for(t=0;t<e.length;t++){var r=e[t],l=r.callback;if(l!==null){if(r.callback=null,r=n,typeof l!="function")throw Error(m(191,l));l.call(r)}}}var Rr={},zt=Jt(Rr),Lr=Jt(Rr),Ir=Jt(Rr);function Cn(e){if(e===Rr)throw Error(m(174));return e}function ri(e,t){switch(pe(Ir,t),pe(Lr,e),pe(zt,Rr),e=t.nodeType,e){case 9:case 11:t=(t=t.documentElement)?t.namespaceURI:sr(null,"");break;default:e=e===8?t.parentNode:t,t=e.namespaceURI||null,e=e.tagName,t=sr(t,e)}ve(zt),pe(zt,t)}function Jn(){ve(zt),ve(Lr),ve(Ir)}function _a(e){Cn(Ir.current);var t=Cn(zt.current),n=sr(t,e.type);t!==n&&(pe(Lr,e),pe(zt,n))}function li(e){Lr.current===e&&(ve(zt),ve(Lr))}var _e=Jt(0);function zl(e){for(var t=e;t!==null;){if(t.tag===13){var n=t.memoizedState;if(n!==null&&(n=n.dehydrated,n===null||n.data==="$?"||n.data==="$!"))return t}else if(t.tag===19&&t.memoizedProps.revealOrder!==void 0){if((t.flags&128)!==0)return t}else if(t.child!==null){t.child.return=t,t=t.child;continue}if(t===e)break;for(;t.sibling===null;){if(t.return===null||t.return===e)return null;t=t.return}t.sibling.return=t.return,t=t.sibling}return null}var oi=[];function ii(){for(var e=0;e<oi.length;e++)oi[e]._workInProgressVersionPrimary=null;oi.length=0}var Rl=Se.ReactCurrentDispatcher,si=Se.ReactCurrentBatchConfig,En=0,we=null,Te=null,Re=null,Ll=!1,Dr=!1,Or=0,pd=0;function Ve(){throw Error(m(321))}function ai(e,t){if(t===null)return!1;for(var n=0;n<t.length&&n<e.length;n++)if(!kt(e[n],t[n]))return!1;return!0}function ui(e,t,n,r,l,o){if(En=o,we=t,t.memoizedState=null,t.updateQueue=null,t.lanes=0,Rl.current=e===null||e.memoizedState===null?gd:yd,e=n(r,l),Dr){o=0;do{if(Dr=!1,Or=0,25<=o)throw Error(m(301));o+=1,Re=Te=null,t.updateQueue=null,Rl.current=_d,e=n(r,l)}while(Dr)}if(Rl.current=Ol,t=Te!==null&&Te.next!==null,En=0,Re=Te=we=null,Ll=!1,t)throw Error(m(300));return e}function ci(){var e=Or!==0;return Or=0,e}function Rt(){var e={memoizedState:null,baseState:null,baseQueue:null,queue:null,next:null};return Re===null?we.memoizedState=Re=e:Re=Re.next=e,Re}function pt(){if(Te===null){var e=we.alternate;e=e!==null?e.memoizedState:null}else e=Te.next;var t=Re===null?we.memoizedState:Re.next;if(t!==null)Re=t,Te=e;else{if(e===null)throw Error(m(310));Te=e,e={memoizedState:Te.memoizedState,baseState:Te.baseState,baseQueue:Te.baseQueue,queue:Te.queue,next:null},Re===null?we.memoizedState=Re=e:Re=Re.next=e}return Re}function Fr(e,t){return typeof t=="function"?t(e):t}function di(e){var t=pt(),n=t.queue;if(n===null)throw Error(m(311));n.lastRenderedReducer=e;var r=Te,l=r.baseQueue,o=n.pending;if(o!==null){if(l!==null){var i=l.next;l.next=o.next,o.next=i}r.baseQueue=l=o,n.pending=null}if(l!==null){o=l.next,r=r.baseState;var s=i=null,d=null,g=o;do{var k=g.lane;if((En&k)===k)d!==null&&(d=d.next={lane:0,action:g.action,hasEagerState:g.hasEagerState,eagerState:g.eagerState,next:null}),r=g.hasEagerState?g.eagerState:e(r,g.action);else{var S={lane:k,action:g.action,hasEagerState:g.hasEagerState,eagerState:g.eagerState,next:null};d===null?(s=d=S,i=r):d=d.next=S,we.lanes|=k,xn|=k}g=g.next}while(g!==null&&g!==o);d===null?i=r:d.next=s,kt(r,t.memoizedState)||(Ze=!0),t.memoizedState=r,t.baseState=i,t.baseQueue=d,n.lastRenderedState=r}if(e=n.interleaved,e!==null){l=e;do o=l.lane,we.lanes|=o,xn|=o,l=l.next;while(l!==e)}else l===null&&(n.lanes=0);return[t.memoizedState,n.dispatch]}function fi(e){var t=pt(),n=t.queue;if(n===null)throw Error(m(311));n.lastRenderedReducer=e;var r=n.dispatch,l=n.pending,o=t.memoizedState;if(l!==null){n.pending=null;var i=l=l.next;do o=e(o,i.action),i=i.next;while(i!==l);kt(o,t.memoizedState)||(Ze=!0),t.memoizedState=o,t.baseQueue===null&&(t.baseState=o),n.lastRenderedState=o}return[o,r]}function wa(){}function ka(e,t){var n=we,r=pt(),l=t(),o=!kt(r.memoizedState,l);if(o&&(r.memoizedState=l,Ze=!0),r=r.queue,pi(Na.bind(null,n,r,e),[e]),r.getSnapshot!==t||o||Re!==null&&Re.memoizedState.tag&1){if(n.flags|=2048,Mr(9,$a.bind(null,n,r,l,t),void 0,null),Le===null)throw Error(m(349));(En&30)!==0||Sa(n,t,l)}return l}function Sa(e,t,n){e.flags|=16384,e={getSnapshot:t,value:n},t=we.updateQueue,t===null?(t={lastEffect:null,stores:null},we.updateQueue=t,t.stores=[e]):(n=t.stores,n===null?t.stores=[e]:n.push(e))}function $a(e,t,n,r){t.value=n,t.getSnapshot=r,Ca(t)&&Ea(e)}function Na(e,t,n){return n(function(){Ca(t)&&Ea(e)})}function Ca(e){var t=e.getSnapshot;e=e.value;try{var n=t();return!kt(e,n)}catch{return!0}}function Ea(e){var t=jt(e,1);t!==null&&Et(t,e,1,-1)}function xa(e){var t=Rt();return typeof e=="function"&&(e=e()),t.memoizedState=t.baseState=e,e={pending:null,interleaved:null,lanes:0,dispatch:null,lastRenderedReducer:Fr,lastRenderedState:e},t.queue=e,e=e.dispatch=vd.bind(null,we,e),[t.memoizedState,e]}function Mr(e,t,n,r){return e={tag:e,create:t,destroy:n,deps:r,next:null},t=we.updateQueue,t===null?(t={lastEffect:null,stores:null},we.updateQueue=t,t.lastEffect=e.next=e):(n=t.lastEffect,n===null?t.lastEffect=e.next=e:(r=n.next,n.next=e,e.next=r,t.lastEffect=e)),e}function Pa(){return pt().memoizedState}function Il(e,t,n,r){var l=Rt();we.flags|=e,l.memoizedState=Mr(1|t,n,void 0,r===void 0?null:r)}function Dl(e,t,n,r){var l=pt();r=r===void 0?null:r;var o=void 0;if(Te!==null){var i=Te.memoizedState;if(o=i.destroy,r!==null&&ai(r,i.deps)){l.memoizedState=Mr(t,n,o,r);return}}we.flags|=e,l.memoizedState=Mr(1|t,n,o,r)}function Ta(e,t){return Il(8390656,8,e,t)}function pi(e,t){return Dl(2048,8,e,t)}function za(e,t){return Dl(4,2,e,t)}function Ra(e,t){return Dl(4,4,e,t)}function La(e,t){if(typeof t=="function")return e=e(),t(e),function(){t(null)};if(t!=null)return e=e(),t.current=e,function(){t.current=null}}function Ia(e,t,n){return n=n!=null?n.concat([e]):null,Dl(4,4,La.bind(null,t,e),n)}function mi(){}function Da(e,t){var n=pt();t=t===void 0?null:t;var r=n.memoizedState;return r!==null&&t!==null&&ai(t,r[1])?r[0]:(n.memoizedState=[e,t],e)}function Oa(e,t){var n=pt();t=t===void 0?null:t;var r=n.memoizedState;return r!==null&&t!==null&&ai(t,r[1])?r[0]:(e=e(),n.memoizedState=[e,t],e)}function Fa(e,t,n){return(En&21)===0?(e.baseState&&(e.baseState=!1,Ze=!0),e.memoizedState=n):(kt(n,t)||(n=fs(),we.lanes|=n,xn|=n,e.baseState=!0),t)}function md(e,t){var n=ae;ae=n!==0&&4>n?n:4,e(!0);var r=si.transition;si.transition={};try{e(!1),t()}finally{ae=n,si.transition=r}}function Ma(){return pt().memoizedState}function hd(e,t,n){var r=sn(e);if(n={lane:r,action:n,hasEagerState:!1,eagerState:null,next:null},ja(e))Aa(t,n);else if(n=ha(e,t,n,r),n!==null){var l=Ge();Et(n,e,r,l),Ua(n,t,r)}}function vd(e,t,n){var r=sn(e),l={lane:r,action:n,hasEagerState:!1,eagerState:null,next:null};if(ja(e))Aa(t,l);else{var o=e.alternate;if(e.lanes===0&&(o===null||o.lanes===0)&&(o=t.lastRenderedReducer,o!==null))try{var i=t.lastRenderedState,s=o(i,n);if(l.hasEagerState=!0,l.eagerState=s,kt(s,i)){var d=t.interleaved;d===null?(l.next=l,ti(t)):(l.next=d.next,d.next=l),t.interleaved=l;return}}catch{}n=ha(e,t,l,r),n!==null&&(l=Ge(),Et(n,e,r,l),Ua(n,t,r))}}function ja(e){var t=e.alternate;return e===we||t!==null&&t===we}function Aa(e,t){Dr=Ll=!0;var n=e.pending;n===null?t.next=t:(t.next=n.next,n.next=t),e.pending=t}function Ua(e,t,n){if((n&4194240)!==0){var r=t.lanes;r&=e.pendingLanes,n|=r,t.lanes=n,vo(e,n)}}var Ol={readContext:ft,useCallback:Ve,useContext:Ve,useEffect:Ve,useImperativeHandle:Ve,useInsertionEffect:Ve,useLayoutEffect:Ve,useMemo:Ve,useReducer:Ve,useRef:Ve,useState:Ve,useDebugValue:Ve,useDeferredValue:Ve,useTransition:Ve,useMutableSource:Ve,useSyncExternalStore:Ve,useId:Ve,unstable_isNewReconciler:!1},gd={readContext:ft,useCallback:function(e,t){return Rt().memoizedState=[e,t===void 0?null:t],e},useContext:ft,useEffect:Ta,useImperativeHandle:function(e,t,n){return n=n!=null?n.concat([e]):null,Il(4194308,4,La.bind(null,t,e),n)},useLayoutEffect:function(e,t){return Il(4194308,4,e,t)},useInsertionEffect:function(e,t){return Il(4,2,e,t)},useMemo:function(e,t){var n=Rt();return t=t===void 0?null:t,e=e(),n.memoizedState=[e,t],e},useReducer:function(e,t,n){var r=Rt();return t=n!==void 0?n(t):t,r.memoizedState=r.baseState=t,e={pending:null,interleaved:null,lanes:0,dispatch:null,lastRenderedReducer:e,lastRenderedState:t},r.queue=e,e=e.dispatch=hd.bind(null,we,e),[r.memoizedState,e]},useRef:function(e){var t=Rt();return e={current:e},t.memoizedState=e},useState:xa,useDebugValue:mi,useDeferredValue:function(e){return Rt().memoizedState=e},useTransition:function(){var e=xa(!1),t=e[0];return e=md.bind(null,e[1]),Rt().memoizedState=e,[t,e]},useMutableSource:function(){},useSyncExternalStore:function(e,t,n){var r=we,l=Rt();if(ge){if(n===void 0)throw Error(m(407));n=n()}else{if(n=t(),Le===null)throw Error(m(349));(En&30)!==0||Sa(r,t,n)}l.memoizedState=n;var o={value:n,getSnapshot:t};return l.queue=o,Ta(Na.bind(null,r,o,e),[e]),r.flags|=2048,Mr(9,$a.bind(null,r,o,n,t),void 0,null),n},useId:function(){var e=Rt(),t=Le.identifierPrefix;if(ge){var n=Mt,r=Ft;n=(r&~(1<<32-wt(r)-1)).toString(32)+n,t=":"+t+"R"+n,n=Or++,0<n&&(t+="H"+n.toString(32)),t+=":"}else n=pd++,t=":"+t+"r"+n.toString(32)+":";return e.memoizedState=t},unstable_isNewReconciler:!1},yd={readContext:ft,useCallback:Da,useContext:ft,useEffect:pi,useImperativeHandle:Ia,useInsertionEffect:za,useLayoutEffect:Ra,useMemo:Oa,useReducer:di,useRef:Pa,useState:function(){return di(Fr)},useDebugValue:mi,useDeferredValue:function(e){var t=pt();return Fa(t,Te.memoizedState,e)},useTransition:function(){var e=di(Fr)[0],t=pt().memoizedState;return[e,t]},useMutableSource:wa,useSyncExternalStore:ka,useId:Ma,unstable_isNewReconciler:!1},_d={readContext:ft,useCallback:Da,useContext:ft,useEffect:pi,useImperativeHandle:Ia,useInsertionEffect:za,useLayoutEffect:Ra,useMemo:Oa,useReducer:fi,useRef:Pa,useState:function(){return fi(Fr)},useDebugValue:mi,useDeferredValue:function(e){var t=pt();return Te===null?t.memoizedState=e:Fa(t,Te.memoizedState,e)},useTransition:function(){var e=fi(Fr)[0],t=pt().memoizedState;return[e,t]},useMutableSource:wa,useSyncExternalStore:ka,useId:Ma,unstable_isNewReconciler:!1};function $t(e,t){if(e&&e.defaultProps){t=$({},t),e=e.defaultProps;for(var n in e)t[n]===void 0&&(t[n]=e[n]);return t}return t}function hi(e,t,n,r){t=e.memoizedState,n=n(r,t),n=n==null?t:$({},t,n),e.memoizedState=n,e.lanes===0&&(e.updateQueue.baseState=n)}var Fl={isMounted:function(e){return(e=e._reactInternals)?_n(e)===e:!1},enqueueSetState:function(e,t,n){e=e._reactInternals;var r=Ge(),l=sn(e),o=At(r,l);o.payload=t,n!=null&&(o.callback=n),t=nn(e,o,l),t!==null&&(Et(t,e,l,r),Pl(t,e,l))},enqueueReplaceState:function(e,t,n){e=e._reactInternals;var r=Ge(),l=sn(e),o=At(r,l);o.tag=1,o.payload=t,n!=null&&(o.callback=n),t=nn(e,o,l),t!==null&&(Et(t,e,l,r),Pl(t,e,l))},enqueueForceUpdate:function(e,t){e=e._reactInternals;var n=Ge(),r=sn(e),l=At(n,r);l.tag=2,t!=null&&(l.callback=t),t=nn(e,l,r),t!==null&&(Et(t,e,r,n),Pl(t,e,r))}};function Ba(e,t,n,r,l,o,i){return e=e.stateNode,typeof e.shouldComponentUpdate=="function"?e.shouldComponentUpdate(r,o,i):t.prototype&&t.prototype.isPureReactComponent?!$r(n,r)||!$r(l,o):!0}function Wa(e,t,n){var r=!1,l=bt,o=t.contextType;return typeof o=="object"&&o!==null?o=ft(o):(l=qe(t)?kn:Qe.current,r=t.contextTypes,o=(r=r!=null)?Hn(e,l):bt),t=new t(n,o),e.memoizedState=t.state!==null&&t.state!==void 0?t.state:null,t.updater=Fl,e.stateNode=t,t._reactInternals=e,r&&(e=e.stateNode,e.__reactInternalMemoizedUnmaskedChildContext=l,e.__reactInternalMemoizedMaskedChildContext=o),t}function Qa(e,t,n,r){e=t.state,typeof t.componentWillReceiveProps=="function"&&t.componentWillReceiveProps(n,r),typeof t.UNSAFE_componentWillReceiveProps=="function"&&t.UNSAFE_componentWillReceiveProps(n,r),t.state!==e&&Fl.enqueueReplaceState(t,t.state,null)}function vi(e,t,n,r){var l=e.stateNode;l.props=n,l.state=e.memoizedState,l.refs={},ni(e);var o=t.contextType;typeof o=="object"&&o!==null?l.context=ft(o):(o=qe(t)?kn:Qe.current,l.context=Hn(e,o)),l.state=e.memoizedState,o=t.getDerivedStateFromProps,typeof o=="function"&&(hi(e,t,o,n),l.state=e.memoizedState),typeof t.getDerivedStateFromProps=="function"||typeof l.getSnapshotBeforeUpdate=="function"||typeof l.UNSAFE_componentWillMount!="function"&&typeof l.componentWillMount!="function"||(t=l.state,typeof l.componentWillMount=="function"&&l.componentWillMount(),typeof l.UNSAFE_componentWillMount=="function"&&l.UNSAFE_componentWillMount(),t!==l.state&&Fl.enqueueReplaceState(l,l.state,null),Tl(e,n,l,r),l.state=e.memoizedState),typeof l.componentDidMount=="function"&&(e.flags|=4194308)}function bn(e,t){try{var n="",r=t;do n+=z(r),r=r.return;while(r);var l=n}catch(o){l=`
Error generating stack: `+o.message+`
`+o.stack}return{value:e,source:t,stack:l,digest:null}}function gi(e,t,n){return{value:e,source:null,stack:n??null,digest:t??null}}function yi(e,t){try{console.error(t.value)}catch(n){setTimeout(function(){throw n})}}var wd=typeof WeakMap=="function"?WeakMap:Map;function Va(e,t,n){n=At(-1,n),n.tag=3,n.payload={element:null};var r=t.value;return n.callback=function(){Ql||(Ql=!0,Ii=r),yi(e,t)},n}function Ha(e,t,n){n=At(-1,n),n.tag=3;var r=e.type.getDerivedStateFromError;if(typeof r=="function"){var l=t.value;n.payload=function(){return r(l)},n.callback=function(){yi(e,t)}}var o=e.stateNode;return o!==null&&typeof o.componentDidCatch=="function"&&(n.callback=function(){yi(e,t),typeof r!="function"&&(ln===null?ln=new Set([this]):ln.add(this));var i=t.stack;this.componentDidCatch(t.value,{componentStack:i!==null?i:""})}),n}function Ka(e,t,n){var r=e.pingCache;if(r===null){r=e.pingCache=new wd;var l=new Set;r.set(t,l)}else l=r.get(t),l===void 0&&(l=new Set,r.set(t,l));l.has(n)||(l.add(n),e=Dd.bind(null,e,t,n),t.then(e,e))}function Ya(e){do{var t;if((t=e.tag===13)&&(t=e.memoizedState,t=t!==null?t.dehydrated!==null:!0),t)return e;e=e.return}while(e!==null);return null}function Ga(e,t,n,r,l){return(e.mode&1)===0?(e===t?e.flags|=65536:(e.flags|=128,n.flags|=131072,n.flags&=-52805,n.tag===1&&(n.alternate===null?n.tag=17:(t=At(-1,1),t.tag=2,nn(n,t,1))),n.lanes|=1),e):(e.flags|=65536,e.lanes=l,e)}var kd=Se.ReactCurrentOwner,Ze=!1;function Ye(e,t,n,r){t.child=e===null?ma(t,null,n,r):Xn(t,e.child,n,r)}function Xa(e,t,n,r,l){n=n.render;var o=t.ref;return Zn(t,l),r=ui(e,t,n,r,o,l),n=ci(),e!==null&&!Ze?(t.updateQueue=e.updateQueue,t.flags&=-2053,e.lanes&=~l,Ut(e,t,l)):(ge&&n&&Ko(t),t.flags|=1,Ye(e,t,r,l),t.child)}function qa(e,t,n,r,l){if(e===null){var o=n.type;return typeof o=="function"&&!Ui(o)&&o.defaultProps===void 0&&n.compare===null&&n.defaultProps===void 0?(t.tag=15,t.type=o,Za(e,t,o,r,l)):(e=Xl(n.type,null,r,t,t.mode,l),e.ref=t.ref,e.return=t,t.child=e)}if(o=e.child,(e.lanes&l)===0){var i=o.memoizedProps;if(n=n.compare,n=n!==null?n:$r,n(i,r)&&e.ref===t.ref)return Ut(e,t,l)}return t.flags|=1,e=un(o,r),e.ref=t.ref,e.return=t,t.child=e}function Za(e,t,n,r,l){if(e!==null){var o=e.memoizedProps;if($r(o,r)&&e.ref===t.ref)if(Ze=!1,t.pendingProps=r=o,(e.lanes&l)!==0)(e.flags&131072)!==0&&(Ze=!0);else return t.lanes=e.lanes,Ut(e,t,l)}return _i(e,t,n,r,l)}function Ja(e,t,n){var r=t.pendingProps,l=r.children,o=e!==null?e.memoizedState:null;if(r.mode==="hidden")if((t.mode&1)===0)t.memoizedState={baseLanes:0,cachePool:null,transitions:null},pe(tr,st),st|=n;else{if((n&1073741824)===0)return e=o!==null?o.baseLanes|n:n,t.lanes=t.childLanes=1073741824,t.memoizedState={baseLanes:e,cachePool:null,transitions:null},t.updateQueue=null,pe(tr,st),st|=e,null;t.memoizedState={baseLanes:0,cachePool:null,transitions:null},r=o!==null?o.baseLanes:n,pe(tr,st),st|=r}else o!==null?(r=o.baseLanes|n,t.memoizedState=null):r=n,pe(tr,st),st|=r;return Ye(e,t,l,n),t.child}function ba(e,t){var n=t.ref;(e===null&&n!==null||e!==null&&e.ref!==n)&&(t.flags|=512,t.flags|=2097152)}function _i(e,t,n,r,l){var o=qe(n)?kn:Qe.current;return o=Hn(t,o),Zn(t,l),n=ui(e,t,n,r,o,l),r=ci(),e!==null&&!Ze?(t.updateQueue=e.updateQueue,t.flags&=-2053,e.lanes&=~l,Ut(e,t,l)):(ge&&r&&Ko(t),t.flags|=1,Ye(e,t,n,l),t.child)}function eu(e,t,n,r,l){if(qe(n)){var o=!0;wl(t)}else o=!1;if(Zn(t,l),t.stateNode===null)jl(e,t),Wa(t,n,r),vi(t,n,r,l),r=!0;else if(e===null){var i=t.stateNode,s=t.memoizedProps;i.props=s;var d=i.context,g=n.contextType;typeof g=="object"&&g!==null?g=ft(g):(g=qe(n)?kn:Qe.current,g=Hn(t,g));var k=n.getDerivedStateFromProps,S=typeof k=="function"||typeof i.getSnapshotBeforeUpdate=="function";S||typeof i.UNSAFE_componentWillReceiveProps!="function"&&typeof i.componentWillReceiveProps!="function"||(s!==r||d!==g)&&Qa(t,i,r,g),tn=!1;var w=t.memoizedState;i.state=w,Tl(t,r,i,l),d=t.memoizedState,s!==r||w!==d||Xe.current||tn?(typeof k=="function"&&(hi(t,n,k,r),d=t.memoizedState),(s=tn||Ba(t,n,s,r,w,d,g))?(S||typeof i.UNSAFE_componentWillMount!="function"&&typeof i.componentWillMount!="function"||(typeof i.componentWillMount=="function"&&i.componentWillMount(),typeof i.UNSAFE_componentWillMount=="function"&&i.UNSAFE_componentWillMount()),typeof i.componentDidMount=="function"&&(t.flags|=4194308)):(typeof i.componentDidMount=="function"&&(t.flags|=4194308),t.memoizedProps=r,t.memoizedState=d),i.props=r,i.state=d,i.context=g,r=s):(typeof i.componentDidMount=="function"&&(t.flags|=4194308),r=!1)}else{i=t.stateNode,va(e,t),s=t.memoizedProps,g=t.type===t.elementType?s:$t(t.type,s),i.props=g,S=t.pendingProps,w=i.context,d=n.contextType,typeof d=="object"&&d!==null?d=ft(d):(d=qe(n)?kn:Qe.current,d=Hn(t,d));var R=n.getDerivedStateFromProps;(k=typeof R=="function"||typeof i.getSnapshotBeforeUpdate=="function")||typeof i.UNSAFE_componentWillReceiveProps!="function"&&typeof i.componentWillReceiveProps!="function"||(s!==S||w!==d)&&Qa(t,i,r,d),tn=!1,w=t.memoizedState,i.state=w,Tl(t,r,i,l);var O=t.memoizedState;s!==S||w!==O||Xe.current||tn?(typeof R=="function"&&(hi(t,n,R,r),O=t.memoizedState),(g=tn||Ba(t,n,g,r,w,O,d)||!1)?(k||typeof i.UNSAFE_componentWillUpdate!="function"&&typeof i.componentWillUpdate!="function"||(typeof i.componentWillUpdate=="function"&&i.componentWillUpdate(r,O,d),typeof i.UNSAFE_componentWillUpdate=="function"&&i.UNSAFE_componentWillUpdate(r,O,d)),typeof i.componentDidUpdate=="function"&&(t.flags|=4),typeof i.getSnapshotBeforeUpdate=="function"&&(t.flags|=1024)):(typeof i.componentDidUpdate!="function"||s===e.memoizedProps&&w===e.memoizedState||(t.flags|=4),typeof i.getSnapshotBeforeUpdate!="function"||s===e.memoizedProps&&w===e.memoizedState||(t.flags|=1024),t.memoizedProps=r,t.memoizedState=O),i.props=r,i.state=O,i.context=d,r=g):(typeof i.componentDidUpdate!="function"||s===e.memoizedProps&&w===e.memoizedState||(t.flags|=4),typeof i.getSnapshotBeforeUpdate!="function"||s===e.memoizedProps&&w===e.memoizedState||(t.flags|=1024),r=!1)}return wi(e,t,n,r,o,l)}function wi(e,t,n,r,l,o){ba(e,t);var i=(t.flags&128)!==0;if(!r&&!i)return l&&oa(t,n,!1),Ut(e,t,o);r=t.stateNode,kd.current=t;var s=i&&typeof n.getDerivedStateFromError!="function"?null:r.render();return t.flags|=1,e!==null&&i?(t.child=Xn(t,e.child,null,o),t.child=Xn(t,null,s,o)):Ye(e,t,s,o),t.memoizedState=r.state,l&&oa(t,n,!0),t.child}function tu(e){var t=e.stateNode;t.pendingContext?ra(e,t.pendingContext,t.pendingContext!==t.context):t.context&&ra(e,t.context,!1),ri(e,t.containerInfo)}function nu(e,t,n,r,l){return Gn(),qo(l),t.flags|=256,Ye(e,t,n,r),t.child}var ki={dehydrated:null,treeContext:null,retryLane:0};function Si(e){return{baseLanes:e,cachePool:null,transitions:null}}function ru(e,t,n){var r=t.pendingProps,l=_e.current,o=!1,i=(t.flags&128)!==0,s;if((s=i)||(s=e!==null&&e.memoizedState===null?!1:(l&2)!==0),s?(o=!0,t.flags&=-129):(e===null||e.memoizedState!==null)&&(l|=1),pe(_e,l&1),e===null)return Xo(t),e=t.memoizedState,e!==null&&(e=e.dehydrated,e!==null)?((t.mode&1)===0?t.lanes=1:e.data==="$!"?t.lanes=8:t.lanes=1073741824,null):(i=r.children,e=r.fallback,o?(r=t.mode,o=t.child,i={mode:"hidden",children:i},(r&1)===0&&o!==null?(o.childLanes=0,o.pendingProps=i):o=ql(i,r,0,null),e=Rn(e,r,n,null),o.return=t,e.return=t,o.sibling=e,t.child=o,t.child.memoizedState=Si(n),t.memoizedState=ki,e):$i(t,i));if(l=e.memoizedState,l!==null&&(s=l.dehydrated,s!==null))return Sd(e,t,i,r,s,l,n);if(o){o=r.fallback,i=t.mode,l=e.child,s=l.sibling;var d={mode:"hidden",children:r.children};return(i&1)===0&&t.child!==l?(r=t.child,r.childLanes=0,r.pendingProps=d,t.deletions=null):(r=un(l,d),r.subtreeFlags=l.subtreeFlags&14680064),s!==null?o=un(s,o):(o=Rn(o,i,n,null),o.flags|=2),o.return=t,r.return=t,r.sibling=o,t.child=r,r=o,o=t.child,i=e.child.memoizedState,i=i===null?Si(n):{baseLanes:i.baseLanes|n,cachePool:null,transitions:i.transitions},o.memoizedState=i,o.childLanes=e.childLanes&~n,t.memoizedState=ki,r}return o=e.child,e=o.sibling,r=un(o,{mode:"visible",children:r.children}),(t.mode&1)===0&&(r.lanes=n),r.return=t,r.sibling=null,e!==null&&(n=t.deletions,n===null?(t.deletions=[e],t.flags|=16):n.push(e)),t.child=r,t.memoizedState=null,r}function $i(e,t){return t=ql({mode:"visible",children:t},e.mode,0,null),t.return=e,e.child=t}function Ml(e,t,n,r){return r!==null&&qo(r),Xn(t,e.child,null,n),e=$i(t,t.pendingProps.children),e.flags|=2,t.memoizedState=null,e}function Sd(e,t,n,r,l,o,i){if(n)return t.flags&256?(t.flags&=-257,r=gi(Error(m(422))),Ml(e,t,i,r)):t.memoizedState!==null?(t.child=e.child,t.flags|=128,null):(o=r.fallback,l=t.mode,r=ql({mode:"visible",children:r.children},l,0,null),o=Rn(o,l,i,null),o.flags|=2,r.return=t,o.return=t,r.sibling=o,t.child=r,(t.mode&1)!==0&&Xn(t,e.child,null,i),t.child.memoizedState=Si(i),t.memoizedState=ki,o);if((t.mode&1)===0)return Ml(e,t,i,null);if(l.data==="$!"){if(r=l.nextSibling&&l.nextSibling.dataset,r)var s=r.dgst;return r=s,o=Error(m(419)),r=gi(o,r,void 0),Ml(e,t,i,r)}if(s=(i&e.childLanes)!==0,Ze||s){if(r=Le,r!==null){switch(i&-i){case 4:l=2;break;case 16:l=8;break;case 64:case 128:case 256:case 512:case 1024:case 2048:case 4096:case 8192:case 16384:case 32768:case 65536:case 131072:case 262144:case 524288:case 1048576:case 2097152:case 4194304:case 8388608:case 16777216:case 33554432:case 67108864:l=32;break;case 536870912:l=268435456;break;default:l=0}l=(l&(r.suspendedLanes|i))!==0?0:l,l!==0&&l!==o.retryLane&&(o.retryLane=l,jt(e,l),Et(r,e,l,-1))}return Ai(),r=gi(Error(m(421))),Ml(e,t,i,r)}return l.data==="$?"?(t.flags|=128,t.child=e.child,t=Od.bind(null,e),l._reactRetry=t,null):(e=o.treeContext,it=Zt(l.nextSibling),ot=t,ge=!0,St=null,e!==null&&(ct[dt++]=Ft,ct[dt++]=Mt,ct[dt++]=Sn,Ft=e.id,Mt=e.overflow,Sn=t),t=$i(t,r.children),t.flags|=4096,t)}function lu(e,t,n){e.lanes|=t;var r=e.alternate;r!==null&&(r.lanes|=t),ei(e.return,t,n)}function Ni(e,t,n,r,l){var o=e.memoizedState;o===null?e.memoizedState={isBackwards:t,rendering:null,renderingStartTime:0,last:r,tail:n,tailMode:l}:(o.isBackwards=t,o.rendering=null,o.renderingStartTime=0,o.last=r,o.tail=n,o.tailMode=l)}function ou(e,t,n){var r=t.pendingProps,l=r.revealOrder,o=r.tail;if(Ye(e,t,r.children,n),r=_e.current,(r&2)!==0)r=r&1|2,t.flags|=128;else{if(e!==null&&(e.flags&128)!==0)e:for(e=t.child;e!==null;){if(e.tag===13)e.memoizedState!==null&&lu(e,n,t);else if(e.tag===19)lu(e,n,t);else if(e.child!==null){e.child.return=e,e=e.child;continue}if(e===t)break e;for(;e.sibling===null;){if(e.return===null||e.return===t)break e;e=e.return}e.sibling.return=e.return,e=e.sibling}r&=1}if(pe(_e,r),(t.mode&1)===0)t.memoizedState=null;else switch(l){case"forwards":for(n=t.child,l=null;n!==null;)e=n.alternate,e!==null&&zl(e)===null&&(l=n),n=n.sibling;n=l,n===null?(l=t.child,t.child=null):(l=n.sibling,n.sibling=null),Ni(t,!1,l,n,o);break;case"backwards":for(n=null,l=t.child,t.child=null;l!==null;){if(e=l.alternate,e!==null&&zl(e)===null){t.child=l;break}e=l.sibling,l.sibling=n,n=l,l=e}Ni(t,!0,n,null,o);break;case"together":Ni(t,!1,null,null,void 0);break;default:t.memoizedState=null}return t.child}function jl(e,t){(t.mode&1)===0&&e!==null&&(e.alternate=null,t.alternate=null,t.flags|=2)}function Ut(e,t,n){if(e!==null&&(t.dependencies=e.dependencies),xn|=t.lanes,(n&t.childLanes)===0)return null;if(e!==null&&t.child!==e.child)throw Error(m(153));if(t.child!==null){for(e=t.child,n=un(e,e.pendingProps),t.child=n,n.return=t;e.sibling!==null;)e=e.sibling,n=n.sibling=un(e,e.pendingProps),n.return=t;n.sibling=null}return t.child}function $d(e,t,n){switch(t.tag){case 3:tu(t),Gn();break;case 5:_a(t);break;case 1:qe(t.type)&&wl(t);break;case 4:ri(t,t.stateNode.containerInfo);break;case 10:var r=t.type._context,l=t.memoizedProps.value;pe(El,r._currentValue),r._currentValue=l;break;case 13:if(r=t.memoizedState,r!==null)return r.dehydrated!==null?(pe(_e,_e.current&1),t.flags|=128,null):(n&t.child.childLanes)!==0?ru(e,t,n):(pe(_e,_e.current&1),e=Ut(e,t,n),e!==null?e.sibling:null);pe(_e,_e.current&1);break;case 19:if(r=(n&t.childLanes)!==0,(e.flags&128)!==0){if(r)return ou(e,t,n);t.flags|=128}if(l=t.memoizedState,l!==null&&(l.rendering=null,l.tail=null,l.lastEffect=null),pe(_e,_e.current),r)break;return null;case 22:case 23:return t.lanes=0,Ja(e,t,n)}return Ut(e,t,n)}var iu,Ci,su,au;iu=function(e,t){for(var n=t.child;n!==null;){if(n.tag===5||n.tag===6)e.appendChild(n.stateNode);else if(n.tag!==4&&n.child!==null){n.child.return=n,n=n.child;continue}if(n===t)break;for(;n.sibling===null;){if(n.return===null||n.return===t)return;n=n.return}n.sibling.return=n.return,n=n.sibling}},Ci=function(){},su=function(e,t,n,r){var l=e.memoizedProps;if(l!==r){e=t.stateNode,Cn(zt.current);var o=null;switch(n){case"input":l=Qt(e,l),r=Qt(e,r),o=[];break;case"select":l=$({},l,{value:void 0}),r=$({},r,{value:void 0}),o=[];break;case"textarea":l=It(e,l),r=It(e,r),o=[];break;default:typeof l.onClick!="function"&&typeof r.onClick=="function"&&(e.onclick=gl)}de(n,r);var i;n=null;for(g in l)if(!r.hasOwnProperty(g)&&l.hasOwnProperty(g)&&l[g]!=null)if(g==="style"){var s=l[g];for(i in s)s.hasOwnProperty(i)&&(n||(n={}),n[i]="")}else g!=="dangerouslySetInnerHTML"&&g!=="children"&&g!=="suppressContentEditableWarning"&&g!=="suppressHydrationWarning"&&g!=="autoFocus"&&(E.hasOwnProperty(g)?o||(o=[]):(o=o||[]).push(g,null));for(g in r){var d=r[g];if(s=l?.[g],r.hasOwnProperty(g)&&d!==s&&(d!=null||s!=null))if(g==="style")if(s){for(i in s)!s.hasOwnProperty(i)||d&&d.hasOwnProperty(i)||(n||(n={}),n[i]="");for(i in d)d.hasOwnProperty(i)&&s[i]!==d[i]&&(n||(n={}),n[i]=d[i])}else n||(o||(o=[]),o.push(g,n)),n=d;else g==="dangerouslySetInnerHTML"?(d=d?d.__html:void 0,s=s?s.__html:void 0,d!=null&&s!==d&&(o=o||[]).push(g,d)):g==="children"?typeof d!="string"&&typeof d!="number"||(o=o||[]).push(g,""+d):g!=="suppressContentEditableWarning"&&g!=="suppressHydrationWarning"&&(E.hasOwnProperty(g)?(d!=null&&g==="onScroll"&&he("scroll",e),o||s===d||(o=[])):(o=o||[]).push(g,d))}n&&(o=o||[]).push("style",n);var g=o;(t.updateQueue=g)&&(t.flags|=4)}},au=function(e,t,n,r){n!==r&&(t.flags|=4)};function jr(e,t){if(!ge)switch(e.tailMode){case"hidden":t=e.tail;for(var n=null;t!==null;)t.alternate!==null&&(n=t),t=t.sibling;n===null?e.tail=null:n.sibling=null;break;case"collapsed":n=e.tail;for(var r=null;n!==null;)n.alternate!==null&&(r=n),n=n.sibling;r===null?t||e.tail===null?e.tail=null:e.tail.sibling=null:r.sibling=null}}function He(e){var t=e.alternate!==null&&e.alternate.child===e.child,n=0,r=0;if(t)for(var l=e.child;l!==null;)n|=l.lanes|l.childLanes,r|=l.subtreeFlags&14680064,r|=l.flags&14680064,l.return=e,l=l.sibling;else for(l=e.child;l!==null;)n|=l.lanes|l.childLanes,r|=l.subtreeFlags,r|=l.flags,l.return=e,l=l.sibling;return e.subtreeFlags|=r,e.childLanes=n,t}function Nd(e,t,n){var r=t.pendingProps;switch(Yo(t),t.tag){case 2:case 16:case 15:case 0:case 11:case 7:case 8:case 12:case 9:case 14:return He(t),null;case 1:return qe(t.type)&&_l(),He(t),null;case 3:return r=t.stateNode,Jn(),ve(Xe),ve(Qe),ii(),r.pendingContext&&(r.context=r.pendingContext,r.pendingContext=null),(e===null||e.child===null)&&(Nl(t)?t.flags|=4:e===null||e.memoizedState.isDehydrated&&(t.flags&256)===0||(t.flags|=1024,St!==null&&(Fi(St),St=null))),Ci(e,t),He(t),null;case 5:li(t);var l=Cn(Ir.current);if(n=t.type,e!==null&&t.stateNode!=null)su(e,t,n,r,l),e.ref!==t.ref&&(t.flags|=512,t.flags|=2097152);else{if(!r){if(t.stateNode===null)throw Error(m(166));return He(t),null}if(e=Cn(zt.current),Nl(t)){r=t.stateNode,n=t.type;var o=t.memoizedProps;switch(r[Tt]=t,r[Pr]=o,e=(t.mode&1)!==0,n){case"dialog":he("cancel",r),he("close",r);break;case"iframe":case"object":case"embed":he("load",r);break;case"video":case"audio":for(l=0;l<Cr.length;l++)he(Cr[l],r);break;case"source":he("error",r);break;case"img":case"image":case"link":he("error",r),he("load",r);break;case"details":he("toggle",r);break;case"input":Vr(r,o),he("invalid",r);break;case"select":r._wrapperState={wasMultiple:!!o.multiple},he("invalid",r);break;case"textarea":Hr(r,o),he("invalid",r)}de(n,o),l=null;for(var i in o)if(o.hasOwnProperty(i)){var s=o[i];i==="children"?typeof s=="string"?r.textContent!==s&&(o.suppressHydrationWarning!==!0&&vl(r.textContent,s,e),l=["children",s]):typeof s=="number"&&r.textContent!==""+s&&(o.suppressHydrationWarning!==!0&&vl(r.textContent,s,e),l=["children",""+s]):E.hasOwnProperty(i)&&s!=null&&i==="onScroll"&&he("scroll",r)}switch(n){case"input":se(r),Ln(r,o,!0);break;case"textarea":se(r),Yr(r);break;case"select":case"option":break;default:typeof o.onClick=="function"&&(r.onclick=gl)}r=l,t.updateQueue=r,r!==null&&(t.flags|=4)}else{i=l.nodeType===9?l:l.ownerDocument,e==="http://www.w3.org/1999/xhtml"&&(e=Gr(n)),e==="http://www.w3.org/1999/xhtml"?n==="script"?(e=i.createElement("div"),e.innerHTML="<script><\/script>",e=e.removeChild(e.firstChild)):typeof r.is=="string"?e=i.createElement(n,{is:r.is}):(e=i.createElement(n),n==="select"&&(i=e,r.multiple?i.multiple=!0:r.size&&(i.size=r.size))):e=i.createElementNS(e,n),e[Tt]=t,e[Pr]=r,iu(e,t,!1,!1),t.stateNode=e;e:{switch(i=_t(n,r),n){case"dialog":he("cancel",e),he("close",e),l=r;break;case"iframe":case"object":case"embed":he("load",e),l=r;break;case"video":case"audio":for(l=0;l<Cr.length;l++)he(Cr[l],e);l=r;break;case"source":he("error",e),l=r;break;case"img":case"image":case"link":he("error",e),he("load",e),l=r;break;case"details":he("toggle",e),l=r;break;case"input":Vr(e,r),l=Qt(e,r),he("invalid",e);break;case"option":l=r;break;case"select":e._wrapperState={wasMultiple:!!r.multiple},l=$({},r,{value:void 0}),he("invalid",e);break;case"textarea":Hr(e,r),l=It(e,r),he("invalid",e);break;default:l=r}de(n,l),s=l;for(o in s)if(s.hasOwnProperty(o)){var d=s[o];o==="style"?L(e,d):o==="dangerouslySetInnerHTML"?(d=d?d.__html:void 0,d!=null&&Xr(e,d)):o==="children"?typeof d=="string"?(n!=="textarea"||d!=="")&&gn(e,d):typeof d=="number"&&gn(e,""+d):o!=="suppressContentEditableWarning"&&o!=="suppressHydrationWarning"&&o!=="autoFocus"&&(E.hasOwnProperty(o)?d!=null&&o==="onScroll"&&he("scroll",e):d!=null&&nt(e,o,d,i))}switch(n){case"input":se(e),Ln(e,r,!1);break;case"textarea":se(e),Yr(e);break;case"option":r.value!=null&&e.setAttribute("value",""+ee(r.value));break;case"select":e.multiple=!!r.multiple,o=r.value,o!=null?Vt(e,!!r.multiple,o,!1):r.defaultValue!=null&&Vt(e,!!r.multiple,r.defaultValue,!0);break;default:typeof l.onClick=="function"&&(e.onclick=gl)}switch(n){case"button":case"input":case"select":case"textarea":r=!!r.autoFocus;break e;case"img":r=!0;break e;default:r=!1}}r&&(t.flags|=4)}t.ref!==null&&(t.flags|=512,t.flags|=2097152)}return He(t),null;case 6:if(e&&t.stateNode!=null)au(e,t,e.memoizedProps,r);else{if(typeof r!="string"&&t.stateNode===null)throw Error(m(166));if(n=Cn(Ir.current),Cn(zt.current),Nl(t)){if(r=t.stateNode,n=t.memoizedProps,r[Tt]=t,(o=r.nodeValue!==n)&&(e=ot,e!==null))switch(e.tag){case 3:vl(r.nodeValue,n,(e.mode&1)!==0);break;case 5:e.memoizedProps.suppressHydrationWarning!==!0&&vl(r.nodeValue,n,(e.mode&1)!==0)}o&&(t.flags|=4)}else r=(n.nodeType===9?n:n.ownerDocument).createTextNode(r),r[Tt]=t,t.stateNode=r}return He(t),null;case 13:if(ve(_e),r=t.memoizedState,e===null||e.memoizedState!==null&&e.memoizedState.dehydrated!==null){if(ge&&it!==null&&(t.mode&1)!==0&&(t.flags&128)===0)da(),Gn(),t.flags|=98560,o=!1;else if(o=Nl(t),r!==null&&r.dehydrated!==null){if(e===null){if(!o)throw Error(m(318));if(o=t.memoizedState,o=o!==null?o.dehydrated:null,!o)throw Error(m(317));o[Tt]=t}else Gn(),(t.flags&128)===0&&(t.memoizedState=null),t.flags|=4;He(t),o=!1}else St!==null&&(Fi(St),St=null),o=!0;if(!o)return t.flags&65536?t:null}return(t.flags&128)!==0?(t.lanes=n,t):(r=r!==null,r!==(e!==null&&e.memoizedState!==null)&&r&&(t.child.flags|=8192,(t.mode&1)!==0&&(e===null||(_e.current&1)!==0?ze===0&&(ze=3):Ai())),t.updateQueue!==null&&(t.flags|=4),He(t),null);case 4:return Jn(),Ci(e,t),e===null&&Er(t.stateNode.containerInfo),He(t),null;case 10:return bo(t.type._context),He(t),null;case 17:return qe(t.type)&&_l(),He(t),null;case 19:if(ve(_e),o=t.memoizedState,o===null)return He(t),null;if(r=(t.flags&128)!==0,i=o.rendering,i===null)if(r)jr(o,!1);else{if(ze!==0||e!==null&&(e.flags&128)!==0)for(e=t.child;e!==null;){if(i=zl(e),i!==null){for(t.flags|=128,jr(o,!1),r=i.updateQueue,r!==null&&(t.updateQueue=r,t.flags|=4),t.subtreeFlags=0,r=n,n=t.child;n!==null;)o=n,e=r,o.flags&=14680066,i=o.alternate,i===null?(o.childLanes=0,o.lanes=e,o.child=null,o.subtreeFlags=0,o.memoizedProps=null,o.memoizedState=null,o.updateQueue=null,o.dependencies=null,o.stateNode=null):(o.childLanes=i.childLanes,o.lanes=i.lanes,o.child=i.child,o.subtreeFlags=0,o.deletions=null,o.memoizedProps=i.memoizedProps,o.memoizedState=i.memoizedState,o.updateQueue=i.updateQueue,o.type=i.type,e=i.dependencies,o.dependencies=e===null?null:{lanes:e.lanes,firstContext:e.firstContext}),n=n.sibling;return pe(_e,_e.current&1|2),t.child}e=e.sibling}o.tail!==null&&$e()>nr&&(t.flags|=128,r=!0,jr(o,!1),t.lanes=4194304)}else{if(!r)if(e=zl(i),e!==null){if(t.flags|=128,r=!0,n=e.updateQueue,n!==null&&(t.updateQueue=n,t.flags|=4),jr(o,!0),o.tail===null&&o.tailMode==="hidden"&&!i.alternate&&!ge)return He(t),null}else 2*$e()-o.renderingStartTime>nr&&n!==1073741824&&(t.flags|=128,r=!0,jr(o,!1),t.lanes=4194304);o.isBackwards?(i.sibling=t.child,t.child=i):(n=o.last,n!==null?n.sibling=i:t.child=i,o.last=i)}return o.tail!==null?(t=o.tail,o.rendering=t,o.tail=t.sibling,o.renderingStartTime=$e(),t.sibling=null,n=_e.current,pe(_e,r?n&1|2:n&1),t):(He(t),null);case 22:case 23:return ji(),r=t.memoizedState!==null,e!==null&&e.memoizedState!==null!==r&&(t.flags|=8192),r&&(t.mode&1)!==0?(st&1073741824)!==0&&(He(t),t.subtreeFlags&6&&(t.flags|=8192)):He(t),null;case 24:return null;case 25:return null}throw Error(m(156,t.tag))}function Cd(e,t){switch(Yo(t),t.tag){case 1:return qe(t.type)&&_l(),e=t.flags,e&65536?(t.flags=e&-65537|128,t):null;case 3:return Jn(),ve(Xe),ve(Qe),ii(),e=t.flags,(e&65536)!==0&&(e&128)===0?(t.flags=e&-65537|128,t):null;case 5:return li(t),null;case 13:if(ve(_e),e=t.memoizedState,e!==null&&e.dehydrated!==null){if(t.alternate===null)throw Error(m(340));Gn()}return e=t.flags,e&65536?(t.flags=e&-65537|128,t):null;case 19:return ve(_e),null;case 4:return Jn(),null;case 10:return bo(t.type._context),null;case 22:case 23:return ji(),null;case 24:return null;default:return null}}var Al=!1,Ke=!1,Ed=typeof WeakSet=="function"?WeakSet:Set,I=null;function er(e,t){var n=e.ref;if(n!==null)if(typeof n=="function")try{n(null)}catch(r){ke(e,t,r)}else n.current=null}function Ei(e,t,n){try{n()}catch(r){ke(e,t,r)}}var uu=!1;function xd(e,t){if(jo=ol,e=Bs(),zo(e)){if("selectionStart"in e)var n={start:e.selectionStart,end:e.selectionEnd};else e:{n=(n=e.ownerDocument)&&n.defaultView||window;var r=n.getSelection&&n.getSelection();if(r&&r.rangeCount!==0){n=r.anchorNode;var l=r.anchorOffset,o=r.focusNode;r=r.focusOffset;try{n.nodeType,o.nodeType}catch{n=null;break e}var i=0,s=-1,d=-1,g=0,k=0,S=e,w=null;t:for(;;){for(var R;S!==n||l!==0&&S.nodeType!==3||(s=i+l),S!==o||r!==0&&S.nodeType!==3||(d=i+r),S.nodeType===3&&(i+=S.nodeValue.length),(R=S.firstChild)!==null;)w=S,S=R;for(;;){if(S===e)break t;if(w===n&&++g===l&&(s=i),w===o&&++k===r&&(d=i),(R=S.nextSibling)!==null)break;S=w,w=S.parentNode}S=R}n=s===-1||d===-1?null:{start:s,end:d}}else n=null}n=n||{start:0,end:0}}else n=null;for(Ao={focusedElem:e,selectionRange:n},ol=!1,I=t;I!==null;)if(t=I,e=t.child,(t.subtreeFlags&1028)!==0&&e!==null)e.return=t,I=e;else for(;I!==null;){t=I;try{var O=t.alternate;if((t.flags&1024)!==0)switch(t.tag){case 0:case 11:case 15:break;case 1:if(O!==null){var F=O.memoizedProps,Ne=O.memoizedState,h=t.stateNode,p=h.getSnapshotBeforeUpdate(t.elementType===t.type?F:$t(t.type,F),Ne);h.__reactInternalSnapshotBeforeUpdate=p}break;case 3:var v=t.stateNode.containerInfo;v.nodeType===1?v.textContent="":v.nodeType===9&&v.documentElement&&v.removeChild(v.documentElement);break;case 5:case 6:case 4:case 17:break;default:throw Error(m(163))}}catch(N){ke(t,t.return,N)}if(e=t.sibling,e!==null){e.return=t.return,I=e;break}I=t.return}return O=uu,uu=!1,O}function Ar(e,t,n){var r=t.updateQueue;if(r=r!==null?r.lastEffect:null,r!==null){var l=r=r.next;do{if((l.tag&e)===e){var o=l.destroy;l.destroy=void 0,o!==void 0&&Ei(t,n,o)}l=l.next}while(l!==r)}}function Ul(e,t){if(t=t.updateQueue,t=t!==null?t.lastEffect:null,t!==null){var n=t=t.next;do{if((n.tag&e)===e){var r=n.create;n.destroy=r()}n=n.next}while(n!==t)}}function xi(e){var t=e.ref;if(t!==null){var n=e.stateNode;e.tag,e=n,typeof t=="function"?t(e):t.current=e}}function cu(e){var t=e.alternate;t!==null&&(e.alternate=null,cu(t)),e.child=null,e.deletions=null,e.sibling=null,e.tag===5&&(t=e.stateNode,t!==null&&(delete t[Tt],delete t[Pr],delete t[Qo],delete t[ud],delete t[cd])),e.stateNode=null,e.return=null,e.dependencies=null,e.memoizedProps=null,e.memoizedState=null,e.pendingProps=null,e.stateNode=null,e.updateQueue=null}function du(e){return e.tag===5||e.tag===3||e.tag===4}function fu(e){e:for(;;){for(;e.sibling===null;){if(e.return===null||du(e.return))return null;e=e.return}for(e.sibling.return=e.return,e=e.sibling;e.tag!==5&&e.tag!==6&&e.tag!==18;){if(e.flags&2||e.child===null||e.tag===4)continue e;e.child.return=e,e=e.child}if(!(e.flags&2))return e.stateNode}}function Pi(e,t,n){var r=e.tag;if(r===5||r===6)e=e.stateNode,t?n.nodeType===8?n.parentNode.insertBefore(e,t):n.insertBefore(e,t):(n.nodeType===8?(t=n.parentNode,t.insertBefore(e,n)):(t=n,t.appendChild(e)),n=n._reactRootContainer,n!=null||t.onclick!==null||(t.onclick=gl));else if(r!==4&&(e=e.child,e!==null))for(Pi(e,t,n),e=e.sibling;e!==null;)Pi(e,t,n),e=e.sibling}function Ti(e,t,n){var r=e.tag;if(r===5||r===6)e=e.stateNode,t?n.insertBefore(e,t):n.appendChild(e);else if(r!==4&&(e=e.child,e!==null))for(Ti(e,t,n),e=e.sibling;e!==null;)Ti(e,t,n),e=e.sibling}var Ae=null,Nt=!1;function rn(e,t,n){for(n=n.child;n!==null;)pu(e,t,n),n=n.sibling}function pu(e,t,n){if(Pt&&typeof Pt.onCommitFiberUnmount=="function")try{Pt.onCommitFiberUnmount(br,n)}catch{}switch(n.tag){case 5:Ke||er(n,t);case 6:var r=Ae,l=Nt;Ae=null,rn(e,t,n),Ae=r,Nt=l,Ae!==null&&(Nt?(e=Ae,n=n.stateNode,e.nodeType===8?e.parentNode.removeChild(n):e.removeChild(n)):Ae.removeChild(n.stateNode));break;case 18:Ae!==null&&(Nt?(e=Ae,n=n.stateNode,e.nodeType===8?Wo(e.parentNode,n):e.nodeType===1&&Wo(e,n),gr(e)):Wo(Ae,n.stateNode));break;case 4:r=Ae,l=Nt,Ae=n.stateNode.containerInfo,Nt=!0,rn(e,t,n),Ae=r,Nt=l;break;case 0:case 11:case 14:case 15:if(!Ke&&(r=n.updateQueue,r!==null&&(r=r.lastEffect,r!==null))){l=r=r.next;do{var o=l,i=o.destroy;o=o.tag,i!==void 0&&((o&2)!==0||(o&4)!==0)&&Ei(n,t,i),l=l.next}while(l!==r)}rn(e,t,n);break;case 1:if(!Ke&&(er(n,t),r=n.stateNode,typeof r.componentWillUnmount=="function"))try{r.props=n.memoizedProps,r.state=n.memoizedState,r.componentWillUnmount()}catch(s){ke(n,t,s)}rn(e,t,n);break;case 21:rn(e,t,n);break;case 22:n.mode&1?(Ke=(r=Ke)||n.memoizedState!==null,rn(e,t,n),Ke=r):rn(e,t,n);break;default:rn(e,t,n)}}function mu(e){var t=e.updateQueue;if(t!==null){e.updateQueue=null;var n=e.stateNode;n===null&&(n=e.stateNode=new Ed),t.forEach(function(r){var l=Fd.bind(null,e,r);n.has(r)||(n.add(r),r.then(l,l))})}}function Ct(e,t){var n=t.deletions;if(n!==null)for(var r=0;r<n.length;r++){var l=n[r];try{var o=e,i=t,s=i;e:for(;s!==null;){switch(s.tag){case 5:Ae=s.stateNode,Nt=!1;break e;case 3:Ae=s.stateNode.containerInfo,Nt=!0;break e;case 4:Ae=s.stateNode.containerInfo,Nt=!0;break e}s=s.return}if(Ae===null)throw Error(m(160));pu(o,i,l),Ae=null,Nt=!1;var d=l.alternate;d!==null&&(d.return=null),l.return=null}catch(g){ke(l,t,g)}}if(t.subtreeFlags&12854)for(t=t.child;t!==null;)hu(t,e),t=t.sibling}function hu(e,t){var n=e.alternate,r=e.flags;switch(e.tag){case 0:case 11:case 14:case 15:if(Ct(t,e),Lt(e),r&4){try{Ar(3,e,e.return),Ul(3,e)}catch(F){ke(e,e.return,F)}try{Ar(5,e,e.return)}catch(F){ke(e,e.return,F)}}break;case 1:Ct(t,e),Lt(e),r&512&&n!==null&&er(n,n.return);break;case 5:if(Ct(t,e),Lt(e),r&512&&n!==null&&er(n,n.return),e.flags&32){var l=e.stateNode;try{gn(l,"")}catch(F){ke(e,e.return,F)}}if(r&4&&(l=e.stateNode,l!=null)){var o=e.memoizedProps,i=n!==null?n.memoizedProps:o,s=e.type,d=e.updateQueue;if(e.updateQueue=null,d!==null)try{s==="input"&&o.type==="radio"&&o.name!=null&&or(l,o),_t(s,i);var g=_t(s,o);for(i=0;i<d.length;i+=2){var k=d[i],S=d[i+1];k==="style"?L(l,S):k==="dangerouslySetInnerHTML"?Xr(l,S):k==="children"?gn(l,S):nt(l,k,S,g)}switch(s){case"input":hn(l,o);break;case"textarea":Kr(l,o);break;case"select":var w=l._wrapperState.wasMultiple;l._wrapperState.wasMultiple=!!o.multiple;var R=o.value;R!=null?Vt(l,!!o.multiple,R,!1):w!==!!o.multiple&&(o.defaultValue!=null?Vt(l,!!o.multiple,o.defaultValue,!0):Vt(l,!!o.multiple,o.multiple?[]:"",!1))}l[Pr]=o}catch(F){ke(e,e.return,F)}}break;case 6:if(Ct(t,e),Lt(e),r&4){if(e.stateNode===null)throw Error(m(162));l=e.stateNode,o=e.memoizedProps;try{l.nodeValue=o}catch(F){ke(e,e.return,F)}}break;case 3:if(Ct(t,e),Lt(e),r&4&&n!==null&&n.memoizedState.isDehydrated)try{gr(t.containerInfo)}catch(F){ke(e,e.return,F)}break;case 4:Ct(t,e),Lt(e);break;case 13:Ct(t,e),Lt(e),l=e.child,l.flags&8192&&(o=l.memoizedState!==null,l.stateNode.isHidden=o,!o||l.alternate!==null&&l.alternate.memoizedState!==null||(Li=$e())),r&4&&mu(e);break;case 22:if(k=n!==null&&n.memoizedState!==null,e.mode&1?(Ke=(g=Ke)||k,Ct(t,e),Ke=g):Ct(t,e),Lt(e),r&8192){if(g=e.memoizedState!==null,(e.stateNode.isHidden=g)&&!k&&(e.mode&1)!==0)for(I=e,k=e.child;k!==null;){for(S=I=k;I!==null;){switch(w=I,R=w.child,w.tag){case 0:case 11:case 14:case 15:Ar(4,w,w.return);break;case 1:er(w,w.return);var O=w.stateNode;if(typeof O.componentWillUnmount=="function"){r=w,n=w.return;try{t=r,O.props=t.memoizedProps,O.state=t.memoizedState,O.componentWillUnmount()}catch(F){ke(r,n,F)}}break;case 5:er(w,w.return);break;case 22:if(w.memoizedState!==null){yu(S);continue}}R!==null?(R.return=w,I=R):yu(S)}k=k.sibling}e:for(k=null,S=e;;){if(S.tag===5){if(k===null){k=S;try{l=S.stateNode,g?(o=l.style,typeof o.setProperty=="function"?o.setProperty("display","none","important"):o.display="none"):(s=S.stateNode,d=S.memoizedProps.style,i=d!=null&&d.hasOwnProperty("display")?d.display:null,s.style.display=u("display",i))}catch(F){ke(e,e.return,F)}}}else if(S.tag===6){if(k===null)try{S.stateNode.nodeValue=g?"":S.memoizedProps}catch(F){ke(e,e.return,F)}}else if((S.tag!==22&&S.tag!==23||S.memoizedState===null||S===e)&&S.child!==null){S.child.return=S,S=S.child;continue}if(S===e)break e;for(;S.sibling===null;){if(S.return===null||S.return===e)break e;k===S&&(k=null),S=S.return}k===S&&(k=null),S.sibling.return=S.return,S=S.sibling}}break;case 19:Ct(t,e),Lt(e),r&4&&mu(e);break;case 21:break;default:Ct(t,e),Lt(e)}}function Lt(e){var t=e.flags;if(t&2){try{e:{for(var n=e.return;n!==null;){if(du(n)){var r=n;break e}n=n.return}throw Error(m(160))}switch(r.tag){case 5:var l=r.stateNode;r.flags&32&&(gn(l,""),r.flags&=-33);var o=fu(e);Ti(e,o,l);break;case 3:case 4:var i=r.stateNode.containerInfo,s=fu(e);Pi(e,s,i);break;default:throw Error(m(161))}}catch(d){ke(e,e.return,d)}e.flags&=-3}t&4096&&(e.flags&=-4097)}function Pd(e,t,n){I=e,vu(e)}function vu(e,t,n){for(var r=(e.mode&1)!==0;I!==null;){var l=I,o=l.child;if(l.tag===22&&r){var i=l.memoizedState!==null||Al;if(!i){var s=l.alternate,d=s!==null&&s.memoizedState!==null||Ke;s=Al;var g=Ke;if(Al=i,(Ke=d)&&!g)for(I=l;I!==null;)i=I,d=i.child,i.tag===22&&i.memoizedState!==null?_u(l):d!==null?(d.return=i,I=d):_u(l);for(;o!==null;)I=o,vu(o),o=o.sibling;I=l,Al=s,Ke=g}gu(e)}else(l.subtreeFlags&8772)!==0&&o!==null?(o.return=l,I=o):gu(e)}}function gu(e){for(;I!==null;){var t=I;if((t.flags&8772)!==0){var n=t.alternate;try{if((t.flags&8772)!==0)switch(t.tag){case 0:case 11:case 15:Ke||Ul(5,t);break;case 1:var r=t.stateNode;if(t.flags&4&&!Ke)if(n===null)r.componentDidMount();else{var l=t.elementType===t.type?n.memoizedProps:$t(t.type,n.memoizedProps);r.componentDidUpdate(l,n.memoizedState,r.__reactInternalSnapshotBeforeUpdate)}var o=t.updateQueue;o!==null&&ya(t,o,r);break;case 3:var i=t.updateQueue;if(i!==null){if(n=null,t.child!==null)switch(t.child.tag){case 5:n=t.child.stateNode;break;case 1:n=t.child.stateNode}ya(t,i,n)}break;case 5:var s=t.stateNode;if(n===null&&t.flags&4){n=s;var d=t.memoizedProps;switch(t.type){case"button":case"input":case"select":case"textarea":d.autoFocus&&n.focus();break;case"img":d.src&&(n.src=d.src)}}break;case 6:break;case 4:break;case 12:break;case 13:if(t.memoizedState===null){var g=t.alternate;if(g!==null){var k=g.memoizedState;if(k!==null){var S=k.dehydrated;S!==null&&gr(S)}}}break;case 19:case 17:case 21:case 22:case 23:case 25:break;default:throw Error(m(163))}Ke||t.flags&512&&xi(t)}catch(w){ke(t,t.return,w)}}if(t===e){I=null;break}if(n=t.sibling,n!==null){n.return=t.return,I=n;break}I=t.return}}function yu(e){for(;I!==null;){var t=I;if(t===e){I=null;break}var n=t.sibling;if(n!==null){n.return=t.return,I=n;break}I=t.return}}function _u(e){for(;I!==null;){var t=I;try{switch(t.tag){case 0:case 11:case 15:var n=t.return;try{Ul(4,t)}catch(d){ke(t,n,d)}break;case 1:var r=t.stateNode;if(typeof r.componentDidMount=="function"){var l=t.return;try{r.componentDidMount()}catch(d){ke(t,l,d)}}var o=t.return;try{xi(t)}catch(d){ke(t,o,d)}break;case 5:var i=t.return;try{xi(t)}catch(d){ke(t,i,d)}}}catch(d){ke(t,t.return,d)}if(t===e){I=null;break}var s=t.sibling;if(s!==null){s.return=t.return,I=s;break}I=t.return}}var Td=Math.ceil,Bl=Se.ReactCurrentDispatcher,zi=Se.ReactCurrentOwner,mt=Se.ReactCurrentBatchConfig,te=0,Le=null,xe=null,Ue=0,st=0,tr=Jt(0),ze=0,Ur=null,xn=0,Wl=0,Ri=0,Br=null,Je=null,Li=0,nr=1/0,Bt=null,Ql=!1,Ii=null,ln=null,Vl=!1,on=null,Hl=0,Wr=0,Di=null,Kl=-1,Yl=0;function Ge(){return(te&6)!==0?$e():Kl!==-1?Kl:Kl=$e()}function sn(e){return(e.mode&1)===0?1:(te&2)!==0&&Ue!==0?Ue&-Ue:fd.transition!==null?(Yl===0&&(Yl=fs()),Yl):(e=ae,e!==0||(e=window.event,e=e===void 0?16:ks(e.type)),e)}function Et(e,t,n,r){if(50<Wr)throw Wr=0,Di=null,Error(m(185));fr(e,n,r),((te&2)===0||e!==Le)&&(e===Le&&((te&2)===0&&(Wl|=n),ze===4&&an(e,Ue)),be(e,r),n===1&&te===0&&(t.mode&1)===0&&(nr=$e()+500,kl&&en()))}function be(e,t){var n=e.callbackNode;fc(e,t);var r=nl(e,e===Le?Ue:0);if(r===0)n!==null&&us(n),e.callbackNode=null,e.callbackPriority=0;else if(t=r&-r,e.callbackPriority!==t){if(n!=null&&us(n),t===1)e.tag===0?dd(ku.bind(null,e)):ia(ku.bind(null,e)),sd(function(){(te&6)===0&&en()}),n=null;else{switch(ps(r)){case 1:n=po;break;case 4:n=cs;break;case 16:n=Jr;break;case 536870912:n=ds;break;default:n=Jr}n=Tu(n,wu.bind(null,e))}e.callbackPriority=t,e.callbackNode=n}}function wu(e,t){if(Kl=-1,Yl=0,(te&6)!==0)throw Error(m(327));var n=e.callbackNode;if(rr()&&e.callbackNode!==n)return null;var r=nl(e,e===Le?Ue:0);if(r===0)return null;if((r&30)!==0||(r&e.expiredLanes)!==0||t)t=Gl(e,r);else{t=r;var l=te;te|=2;var o=$u();(Le!==e||Ue!==t)&&(Bt=null,nr=$e()+500,Tn(e,t));do try{Ld();break}catch(s){Su(e,s)}while(!0);Jo(),Bl.current=o,te=l,xe!==null?t=0:(Le=null,Ue=0,t=ze)}if(t!==0){if(t===2&&(l=mo(e),l!==0&&(r=l,t=Oi(e,l))),t===1)throw n=Ur,Tn(e,0),an(e,r),be(e,$e()),n;if(t===6)an(e,r);else{if(l=e.current.alternate,(r&30)===0&&!zd(l)&&(t=Gl(e,r),t===2&&(o=mo(e),o!==0&&(r=o,t=Oi(e,o))),t===1))throw n=Ur,Tn(e,0),an(e,r),be(e,$e()),n;switch(e.finishedWork=l,e.finishedLanes=r,t){case 0:case 1:throw Error(m(345));case 2:zn(e,Je,Bt);break;case 3:if(an(e,r),(r&130023424)===r&&(t=Li+500-$e(),10<t)){if(nl(e,0)!==0)break;if(l=e.suspendedLanes,(l&r)!==r){Ge(),e.pingedLanes|=e.suspendedLanes&l;break}e.timeoutHandle=Bo(zn.bind(null,e,Je,Bt),t);break}zn(e,Je,Bt);break;case 4:if(an(e,r),(r&4194240)===r)break;for(t=e.eventTimes,l=-1;0<r;){var i=31-wt(r);o=1<<i,i=t[i],i>l&&(l=i),r&=~o}if(r=l,r=$e()-r,r=(120>r?120:480>r?480:1080>r?1080:1920>r?1920:3e3>r?3e3:4320>r?4320:1960*Td(r/1960))-r,10<r){e.timeoutHandle=Bo(zn.bind(null,e,Je,Bt),r);break}zn(e,Je,Bt);break;case 5:zn(e,Je,Bt);break;default:throw Error(m(329))}}}return be(e,$e()),e.callbackNode===n?wu.bind(null,e):null}function Oi(e,t){var n=Br;return e.current.memoizedState.isDehydrated&&(Tn(e,t).flags|=256),e=Gl(e,t),e!==2&&(t=Je,Je=n,t!==null&&Fi(t)),e}function Fi(e){Je===null?Je=e:Je.push.apply(Je,e)}function zd(e){for(var t=e;;){if(t.flags&16384){var n=t.updateQueue;if(n!==null&&(n=n.stores,n!==null))for(var r=0;r<n.length;r++){var l=n[r],o=l.getSnapshot;l=l.value;try{if(!kt(o(),l))return!1}catch{return!1}}}if(n=t.child,t.subtreeFlags&16384&&n!==null)n.return=t,t=n;else{if(t===e)break;for(;t.sibling===null;){if(t.return===null||t.return===e)return!0;t=t.return}t.sibling.return=t.return,t=t.sibling}}return!0}function an(e,t){for(t&=~Ri,t&=~Wl,e.suspendedLanes|=t,e.pingedLanes&=~t,e=e.expirationTimes;0<t;){var n=31-wt(t),r=1<<n;e[n]=-1,t&=~r}}function ku(e){if((te&6)!==0)throw Error(m(327));rr();var t=nl(e,0);if((t&1)===0)return be(e,$e()),null;var n=Gl(e,t);if(e.tag!==0&&n===2){var r=mo(e);r!==0&&(t=r,n=Oi(e,r))}if(n===1)throw n=Ur,Tn(e,0),an(e,t),be(e,$e()),n;if(n===6)throw Error(m(345));return e.finishedWork=e.current.alternate,e.finishedLanes=t,zn(e,Je,Bt),be(e,$e()),null}function Mi(e,t){var n=te;te|=1;try{return e(t)}finally{te=n,te===0&&(nr=$e()+500,kl&&en())}}function Pn(e){on!==null&&on.tag===0&&(te&6)===0&&rr();var t=te;te|=1;var n=mt.transition,r=ae;try{if(mt.transition=null,ae=1,e)return e()}finally{ae=r,mt.transition=n,te=t,(te&6)===0&&en()}}function ji(){st=tr.current,ve(tr)}function Tn(e,t){e.finishedWork=null,e.finishedLanes=0;var n=e.timeoutHandle;if(n!==-1&&(e.timeoutHandle=-1,id(n)),xe!==null)for(n=xe.return;n!==null;){var r=n;switch(Yo(r),r.tag){case 1:r=r.type.childContextTypes,r!=null&&_l();break;case 3:Jn(),ve(Xe),ve(Qe),ii();break;case 5:li(r);break;case 4:Jn();break;case 13:ve(_e);break;case 19:ve(_e);break;case 10:bo(r.type._context);break;case 22:case 23:ji()}n=n.return}if(Le=e,xe=e=un(e.current,null),Ue=st=t,ze=0,Ur=null,Ri=Wl=xn=0,Je=Br=null,Nn!==null){for(t=0;t<Nn.length;t++)if(n=Nn[t],r=n.interleaved,r!==null){n.interleaved=null;var l=r.next,o=n.pending;if(o!==null){var i=o.next;o.next=l,r.next=i}n.pending=r}Nn=null}return e}function Su(e,t){do{var n=xe;try{if(Jo(),Rl.current=Ol,Ll){for(var r=we.memoizedState;r!==null;){var l=r.queue;l!==null&&(l.pending=null),r=r.next}Ll=!1}if(En=0,Re=Te=we=null,Dr=!1,Or=0,zi.current=null,n===null||n.return===null){ze=1,Ur=t,xe=null;break}e:{var o=e,i=n.return,s=n,d=t;if(t=Ue,s.flags|=32768,d!==null&&typeof d=="object"&&typeof d.then=="function"){var g=d,k=s,S=k.tag;if((k.mode&1)===0&&(S===0||S===11||S===15)){var w=k.alternate;w?(k.updateQueue=w.updateQueue,k.memoizedState=w.memoizedState,k.lanes=w.lanes):(k.updateQueue=null,k.memoizedState=null)}var R=Ya(i);if(R!==null){R.flags&=-257,Ga(R,i,s,o,t),R.mode&1&&Ka(o,g,t),t=R,d=g;var O=t.updateQueue;if(O===null){var F=new Set;F.add(d),t.updateQueue=F}else O.add(d);break e}else{if((t&1)===0){Ka(o,g,t),Ai();break e}d=Error(m(426))}}else if(ge&&s.mode&1){var Ne=Ya(i);if(Ne!==null){(Ne.flags&65536)===0&&(Ne.flags|=256),Ga(Ne,i,s,o,t),qo(bn(d,s));break e}}o=d=bn(d,s),ze!==4&&(ze=2),Br===null?Br=[o]:Br.push(o),o=i;do{switch(o.tag){case 3:o.flags|=65536,t&=-t,o.lanes|=t;var h=Va(o,d,t);ga(o,h);break e;case 1:s=d;var p=o.type,v=o.stateNode;if((o.flags&128)===0&&(typeof p.getDerivedStateFromError=="function"||v!==null&&typeof v.componentDidCatch=="function"&&(ln===null||!ln.has(v)))){o.flags|=65536,t&=-t,o.lanes|=t;var N=Ha(o,s,t);ga(o,N);break e}}o=o.return}while(o!==null)}Cu(n)}catch(A){t=A,xe===n&&n!==null&&(xe=n=n.return);continue}break}while(!0)}function $u(){var e=Bl.current;return Bl.current=Ol,e===null?Ol:e}function Ai(){(ze===0||ze===3||ze===2)&&(ze=4),Le===null||(xn&268435455)===0&&(Wl&268435455)===0||an(Le,Ue)}function Gl(e,t){var n=te;te|=2;var r=$u();(Le!==e||Ue!==t)&&(Bt=null,Tn(e,t));do try{Rd();break}catch(l){Su(e,l)}while(!0);if(Jo(),te=n,Bl.current=r,xe!==null)throw Error(m(261));return Le=null,Ue=0,ze}function Rd(){for(;xe!==null;)Nu(xe)}function Ld(){for(;xe!==null&&!rc();)Nu(xe)}function Nu(e){var t=Pu(e.alternate,e,st);e.memoizedProps=e.pendingProps,t===null?Cu(e):xe=t,zi.current=null}function Cu(e){var t=e;do{var n=t.alternate;if(e=t.return,(t.flags&32768)===0){if(n=Nd(n,t,st),n!==null){xe=n;return}}else{if(n=Cd(n,t),n!==null){n.flags&=32767,xe=n;return}if(e!==null)e.flags|=32768,e.subtreeFlags=0,e.deletions=null;else{ze=6,xe=null;return}}if(t=t.sibling,t!==null){xe=t;return}xe=t=e}while(t!==null);ze===0&&(ze=5)}function zn(e,t,n){var r=ae,l=mt.transition;try{mt.transition=null,ae=1,Id(e,t,n,r)}finally{mt.transition=l,ae=r}return null}function Id(e,t,n,r){do rr();while(on!==null);if((te&6)!==0)throw Error(m(327));n=e.finishedWork;var l=e.finishedLanes;if(n===null)return null;if(e.finishedWork=null,e.finishedLanes=0,n===e.current)throw Error(m(177));e.callbackNode=null,e.callbackPriority=0;var o=n.lanes|n.childLanes;if(pc(e,o),e===Le&&(xe=Le=null,Ue=0),(n.subtreeFlags&2064)===0&&(n.flags&2064)===0||Vl||(Vl=!0,Tu(Jr,function(){return rr(),null})),o=(n.flags&15990)!==0,(n.subtreeFlags&15990)!==0||o){o=mt.transition,mt.transition=null;var i=ae;ae=1;var s=te;te|=4,zi.current=null,xd(e,n),hu(n,e),bc(Ao),ol=!!jo,Ao=jo=null,e.current=n,Pd(n),lc(),te=s,ae=i,mt.transition=o}else e.current=n;if(Vl&&(Vl=!1,on=e,Hl=l),o=e.pendingLanes,o===0&&(ln=null),sc(n.stateNode),be(e,$e()),t!==null)for(r=e.onRecoverableError,n=0;n<t.length;n++)l=t[n],r(l.value,{componentStack:l.stack,digest:l.digest});if(Ql)throw Ql=!1,e=Ii,Ii=null,e;return(Hl&1)!==0&&e.tag!==0&&rr(),o=e.pendingLanes,(o&1)!==0?e===Di?Wr++:(Wr=0,Di=e):Wr=0,en(),null}function rr(){if(on!==null){var e=ps(Hl),t=mt.transition,n=ae;try{if(mt.transition=null,ae=16>e?16:e,on===null)var r=!1;else{if(e=on,on=null,Hl=0,(te&6)!==0)throw Error(m(331));var l=te;for(te|=4,I=e.current;I!==null;){var o=I,i=o.child;if((I.flags&16)!==0){var s=o.deletions;if(s!==null){for(var d=0;d<s.length;d++){var g=s[d];for(I=g;I!==null;){var k=I;switch(k.tag){case 0:case 11:case 15:Ar(8,k,o)}var S=k.child;if(S!==null)S.return=k,I=S;else for(;I!==null;){k=I;var w=k.sibling,R=k.return;if(cu(k),k===g){I=null;break}if(w!==null){w.return=R,I=w;break}I=R}}}var O=o.alternate;if(O!==null){var F=O.child;if(F!==null){O.child=null;do{var Ne=F.sibling;F.sibling=null,F=Ne}while(F!==null)}}I=o}}if((o.subtreeFlags&2064)!==0&&i!==null)i.return=o,I=i;else e:for(;I!==null;){if(o=I,(o.flags&2048)!==0)switch(o.tag){case 0:case 11:case 15:Ar(9,o,o.return)}var h=o.sibling;if(h!==null){h.return=o.return,I=h;break e}I=o.return}}var p=e.current;for(I=p;I!==null;){i=I;var v=i.child;if((i.subtreeFlags&2064)!==0&&v!==null)v.return=i,I=v;else e:for(i=p;I!==null;){if(s=I,(s.flags&2048)!==0)try{switch(s.tag){case 0:case 11:case 15:Ul(9,s)}}catch(A){ke(s,s.return,A)}if(s===i){I=null;break e}var N=s.sibling;if(N!==null){N.return=s.return,I=N;break e}I=s.return}}if(te=l,en(),Pt&&typeof Pt.onPostCommitFiberRoot=="function")try{Pt.onPostCommitFiberRoot(br,e)}catch{}r=!0}return r}finally{ae=n,mt.transition=t}}return!1}function Eu(e,t,n){t=bn(n,t),t=Va(e,t,1),e=nn(e,t,1),t=Ge(),e!==null&&(fr(e,1,t),be(e,t))}function ke(e,t,n){if(e.tag===3)Eu(e,e,n);else for(;t!==null;){if(t.tag===3){Eu(t,e,n);break}else if(t.tag===1){var r=t.stateNode;if(typeof t.type.getDerivedStateFromError=="function"||typeof r.componentDidCatch=="function"&&(ln===null||!ln.has(r))){e=bn(n,e),e=Ha(t,e,1),t=nn(t,e,1),e=Ge(),t!==null&&(fr(t,1,e),be(t,e));break}}t=t.return}}function Dd(e,t,n){var r=e.pingCache;r!==null&&r.delete(t),t=Ge(),e.pingedLanes|=e.suspendedLanes&n,Le===e&&(Ue&n)===n&&(ze===4||ze===3&&(Ue&130023424)===Ue&&500>$e()-Li?Tn(e,0):Ri|=n),be(e,t)}function xu(e,t){t===0&&((e.mode&1)===0?t=1:(t=tl,tl<<=1,(tl&130023424)===0&&(tl=4194304)));var n=Ge();e=jt(e,t),e!==null&&(fr(e,t,n),be(e,n))}function Od(e){var t=e.memoizedState,n=0;t!==null&&(n=t.retryLane),xu(e,n)}function Fd(e,t){var n=0;switch(e.tag){case 13:var r=e.stateNode,l=e.memoizedState;l!==null&&(n=l.retryLane);break;case 19:r=e.stateNode;break;default:throw Error(m(314))}r!==null&&r.delete(t),xu(e,n)}var Pu;Pu=function(e,t,n){if(e!==null)if(e.memoizedProps!==t.pendingProps||Xe.current)Ze=!0;else{if((e.lanes&n)===0&&(t.flags&128)===0)return Ze=!1,$d(e,t,n);Ze=(e.flags&131072)!==0}else Ze=!1,ge&&(t.flags&1048576)!==0&&sa(t,$l,t.index);switch(t.lanes=0,t.tag){case 2:var r=t.type;jl(e,t),e=t.pendingProps;var l=Hn(t,Qe.current);Zn(t,n),l=ui(null,t,r,e,l,n);var o=ci();return t.flags|=1,typeof l=="object"&&l!==null&&typeof l.render=="function"&&l.$$typeof===void 0?(t.tag=1,t.memoizedState=null,t.updateQueue=null,qe(r)?(o=!0,wl(t)):o=!1,t.memoizedState=l.state!==null&&l.state!==void 0?l.state:null,ni(t),l.updater=Fl,t.stateNode=l,l._reactInternals=t,vi(t,r,e,n),t=wi(null,t,r,!0,o,n)):(t.tag=0,ge&&o&&Ko(t),Ye(null,t,l,n),t=t.child),t;case 16:r=t.elementType;e:{switch(jl(e,t),e=t.pendingProps,l=r._init,r=l(r._payload),t.type=r,l=t.tag=jd(r),e=$t(r,e),l){case 0:t=_i(null,t,r,e,n);break e;case 1:t=eu(null,t,r,e,n);break e;case 11:t=Xa(null,t,r,e,n);break e;case 14:t=qa(null,t,r,$t(r.type,e),n);break e}throw Error(m(306,r,""))}return t;case 0:return r=t.type,l=t.pendingProps,l=t.elementType===r?l:$t(r,l),_i(e,t,r,l,n);case 1:return r=t.type,l=t.pendingProps,l=t.elementType===r?l:$t(r,l),eu(e,t,r,l,n);case 3:e:{if(tu(t),e===null)throw Error(m(387));r=t.pendingProps,o=t.memoizedState,l=o.element,va(e,t),Tl(t,r,null,n);var i=t.memoizedState;if(r=i.element,o.isDehydrated)if(o={element:r,isDehydrated:!1,cache:i.cache,pendingSuspenseBoundaries:i.pendingSuspenseBoundaries,transitions:i.transitions},t.updateQueue.baseState=o,t.memoizedState=o,t.flags&256){l=bn(Error(m(423)),t),t=nu(e,t,r,n,l);break e}else if(r!==l){l=bn(Error(m(424)),t),t=nu(e,t,r,n,l);break e}else for(it=Zt(t.stateNode.containerInfo.firstChild),ot=t,ge=!0,St=null,n=ma(t,null,r,n),t.child=n;n;)n.flags=n.flags&-3|4096,n=n.sibling;else{if(Gn(),r===l){t=Ut(e,t,n);break e}Ye(e,t,r,n)}t=t.child}return t;case 5:return _a(t),e===null&&Xo(t),r=t.type,l=t.pendingProps,o=e!==null?e.memoizedProps:null,i=l.children,Uo(r,l)?i=null:o!==null&&Uo(r,o)&&(t.flags|=32),ba(e,t),Ye(e,t,i,n),t.child;case 6:return e===null&&Xo(t),null;case 13:return ru(e,t,n);case 4:return ri(t,t.stateNode.containerInfo),r=t.pendingProps,e===null?t.child=Xn(t,null,r,n):Ye(e,t,r,n),t.child;case 11:return r=t.type,l=t.pendingProps,l=t.elementType===r?l:$t(r,l),Xa(e,t,r,l,n);case 7:return Ye(e,t,t.pendingProps,n),t.child;case 8:return Ye(e,t,t.pendingProps.children,n),t.child;case 12:return Ye(e,t,t.pendingProps.children,n),t.child;case 10:e:{if(r=t.type._context,l=t.pendingProps,o=t.memoizedProps,i=l.value,pe(El,r._currentValue),r._currentValue=i,o!==null)if(kt(o.value,i)){if(o.children===l.children&&!Xe.current){t=Ut(e,t,n);break e}}else for(o=t.child,o!==null&&(o.return=t);o!==null;){var s=o.dependencies;if(s!==null){i=o.child;for(var d=s.firstContext;d!==null;){if(d.context===r){if(o.tag===1){d=At(-1,n&-n),d.tag=2;var g=o.updateQueue;if(g!==null){g=g.shared;var k=g.pending;k===null?d.next=d:(d.next=k.next,k.next=d),g.pending=d}}o.lanes|=n,d=o.alternate,d!==null&&(d.lanes|=n),ei(o.return,n,t),s.lanes|=n;break}d=d.next}}else if(o.tag===10)i=o.type===t.type?null:o.child;else if(o.tag===18){if(i=o.return,i===null)throw Error(m(341));i.lanes|=n,s=i.alternate,s!==null&&(s.lanes|=n),ei(i,n,t),i=o.sibling}else i=o.child;if(i!==null)i.return=o;else for(i=o;i!==null;){if(i===t){i=null;break}if(o=i.sibling,o!==null){o.return=i.return,i=o;break}i=i.return}o=i}Ye(e,t,l.children,n),t=t.child}return t;case 9:return l=t.type,r=t.pendingProps.children,Zn(t,n),l=ft(l),r=r(l),t.flags|=1,Ye(e,t,r,n),t.child;case 14:return r=t.type,l=$t(r,t.pendingProps),l=$t(r.type,l),qa(e,t,r,l,n);case 15:return Za(e,t,t.type,t.pendingProps,n);case 17:return r=t.type,l=t.pendingProps,l=t.elementType===r?l:$t(r,l),jl(e,t),t.tag=1,qe(r)?(e=!0,wl(t)):e=!1,Zn(t,n),Wa(t,r,l),vi(t,r,l,n),wi(null,t,r,!0,e,n);case 19:return ou(e,t,n);case 22:return Ja(e,t,n)}throw Error(m(156,t.tag))};function Tu(e,t){return as(e,t)}function Md(e,t,n,r){this.tag=e,this.key=n,this.sibling=this.child=this.return=this.stateNode=this.type=this.elementType=null,this.index=0,this.ref=null,this.pendingProps=t,this.dependencies=this.memoizedState=this.updateQueue=this.memoizedProps=null,this.mode=r,this.subtreeFlags=this.flags=0,this.deletions=null,this.childLanes=this.lanes=0,this.alternate=null}function ht(e,t,n,r){return new Md(e,t,n,r)}function Ui(e){return e=e.prototype,!(!e||!e.isReactComponent)}function jd(e){if(typeof e=="function")return Ui(e)?1:0;if(e!=null){if(e=e.$$typeof,e===J)return 11;if(e===We)return 14}return 2}function un(e,t){var n=e.alternate;return n===null?(n=ht(e.tag,t,e.key,e.mode),n.elementType=e.elementType,n.type=e.type,n.stateNode=e.stateNode,n.alternate=e,e.alternate=n):(n.pendingProps=t,n.type=e.type,n.flags=0,n.subtreeFlags=0,n.deletions=null),n.flags=e.flags&14680064,n.childLanes=e.childLanes,n.lanes=e.lanes,n.child=e.child,n.memoizedProps=e.memoizedProps,n.memoizedState=e.memoizedState,n.updateQueue=e.updateQueue,t=e.dependencies,n.dependencies=t===null?null:{lanes:t.lanes,firstContext:t.firstContext},n.sibling=e.sibling,n.index=e.index,n.ref=e.ref,n}function Xl(e,t,n,r,l,o){var i=2;if(r=e,typeof e=="function")Ui(e)&&(i=1);else if(typeof e=="string")i=5;else e:switch(e){case Me:return Rn(n.children,l,o,t);case Be:i=8,l|=8;break;case xt:return e=ht(12,n,t,l|2),e.elementType=xt,e.lanes=o,e;case me:return e=ht(13,n,t,l),e.elementType=me,e.lanes=o,e;case Pe:return e=ht(19,n,t,l),e.elementType=Pe,e.lanes=o,e;case Q:return ql(n,l,o,t);default:if(typeof e=="object"&&e!==null)switch(e.$$typeof){case gt:i=10;break e;case yt:i=9;break e;case J:i=11;break e;case We:i=14;break e;case j:i=16,r=null;break e}throw Error(m(130,e==null?e:typeof e,""))}return t=ht(i,n,t,l),t.elementType=e,t.type=r,t.lanes=o,t}function Rn(e,t,n,r){return e=ht(7,e,r,t),e.lanes=n,e}function ql(e,t,n,r){return e=ht(22,e,r,t),e.elementType=Q,e.lanes=n,e.stateNode={isHidden:!1},e}function Bi(e,t,n){return e=ht(6,e,null,t),e.lanes=n,e}function Wi(e,t,n){return t=ht(4,e.children!==null?e.children:[],e.key,t),t.lanes=n,t.stateNode={containerInfo:e.containerInfo,pendingChildren:null,implementation:e.implementation},t}function Ad(e,t,n,r,l){this.tag=t,this.containerInfo=e,this.finishedWork=this.pingCache=this.current=this.pendingChildren=null,this.timeoutHandle=-1,this.callbackNode=this.pendingContext=this.context=null,this.callbackPriority=0,this.eventTimes=ho(0),this.expirationTimes=ho(-1),this.entangledLanes=this.finishedLanes=this.mutableReadLanes=this.expiredLanes=this.pingedLanes=this.suspendedLanes=this.pendingLanes=0,this.entanglements=ho(0),this.identifierPrefix=r,this.onRecoverableError=l,this.mutableSourceEagerHydrationData=null}function Qi(e,t,n,r,l,o,i,s,d){return e=new Ad(e,t,n,s,d),t===1?(t=1,o===!0&&(t|=8)):t=0,o=ht(3,null,null,t),e.current=o,o.stateNode=e,o.memoizedState={element:r,isDehydrated:n,cache:null,transitions:null,pendingSuspenseBoundaries:null},ni(o),e}function Ud(e,t,n){var r=3<arguments.length&&arguments[3]!==void 0?arguments[3]:null;return{$$typeof:Fe,key:r==null?null:""+r,children:e,containerInfo:t,implementation:n}}function zu(e){if(!e)return bt;e=e._reactInternals;e:{if(_n(e)!==e||e.tag!==1)throw Error(m(170));var t=e;do{switch(t.tag){case 3:t=t.stateNode.context;break e;case 1:if(qe(t.type)){t=t.stateNode.__reactInternalMemoizedMergedChildContext;break e}}t=t.return}while(t!==null);throw Error(m(171))}if(e.tag===1){var n=e.type;if(qe(n))return la(e,n,t)}return t}function Ru(e,t,n,r,l,o,i,s,d){return e=Qi(n,r,!0,e,l,o,i,s,d),e.context=zu(null),n=e.current,r=Ge(),l=sn(n),o=At(r,l),o.callback=t??null,nn(n,o,l),e.current.lanes=l,fr(e,l,r),be(e,r),e}function Zl(e,t,n,r){var l=t.current,o=Ge(),i=sn(l);return n=zu(n),t.context===null?t.context=n:t.pendingContext=n,t=At(o,i),t.payload={element:e},r=r===void 0?null:r,r!==null&&(t.callback=r),e=nn(l,t,i),e!==null&&(Et(e,l,i,o),Pl(e,l,i)),i}function Jl(e){return e=e.current,e.child?(e.child.tag===5,e.child.stateNode):null}function Lu(e,t){if(e=e.memoizedState,e!==null&&e.dehydrated!==null){var n=e.retryLane;e.retryLane=n!==0&&n<t?n:t}}function Vi(e,t){Lu(e,t),(e=e.alternate)&&Lu(e,t)}function Bd(){return null}var Iu=typeof reportError=="function"?reportError:function(e){console.error(e)};function Hi(e){this._internalRoot=e}bl.prototype.render=Hi.prototype.render=function(e){var t=this._internalRoot;if(t===null)throw Error(m(409));Zl(e,t,null,null)},bl.prototype.unmount=Hi.prototype.unmount=function(){var e=this._internalRoot;if(e!==null){this._internalRoot=null;var t=e.containerInfo;Pn(function(){Zl(null,e,null,null)}),t[Dt]=null}};function bl(e){this._internalRoot=e}bl.prototype.unstable_scheduleHydration=function(e){if(e){var t=vs();e={blockedOn:null,target:e,priority:t};for(var n=0;n<Gt.length&&t!==0&&t<Gt[n].priority;n++);Gt.splice(n,0,e),n===0&&_s(e)}};function Ki(e){return!(!e||e.nodeType!==1&&e.nodeType!==9&&e.nodeType!==11)}function eo(e){return!(!e||e.nodeType!==1&&e.nodeType!==9&&e.nodeType!==11&&(e.nodeType!==8||e.nodeValue!==" react-mount-point-unstable "))}function Du(){}function Wd(e,t,n,r,l){if(l){if(typeof r=="function"){var o=r;r=function(){var g=Jl(i);o.call(g)}}var i=Ru(t,r,e,0,null,!1,!1,"",Du);return e._reactRootContainer=i,e[Dt]=i.current,Er(e.nodeType===8?e.parentNode:e),Pn(),i}for(;l=e.lastChild;)e.removeChild(l);if(typeof r=="function"){var s=r;r=function(){var g=Jl(d);s.call(g)}}var d=Qi(e,0,!1,null,null,!1,!1,"",Du);return e._reactRootContainer=d,e[Dt]=d.current,Er(e.nodeType===8?e.parentNode:e),Pn(function(){Zl(t,d,n,r)}),d}function to(e,t,n,r,l){var o=n._reactRootContainer;if(o){var i=o;if(typeof l=="function"){var s=l;l=function(){var d=Jl(i);s.call(d)}}Zl(t,i,e,l)}else i=Wd(n,t,e,l,r);return Jl(i)}ms=function(e){switch(e.tag){case 3:var t=e.stateNode;if(t.current.memoizedState.isDehydrated){var n=dr(t.pendingLanes);n!==0&&(vo(t,n|1),be(t,$e()),(te&6)===0&&(nr=$e()+500,en()))}break;case 13:Pn(function(){var r=jt(e,1);if(r!==null){var l=Ge();Et(r,e,1,l)}}),Vi(e,1)}},go=function(e){if(e.tag===13){var t=jt(e,134217728);if(t!==null){var n=Ge();Et(t,e,134217728,n)}Vi(e,134217728)}},hs=function(e){if(e.tag===13){var t=sn(e),n=jt(e,t);if(n!==null){var r=Ge();Et(n,e,t,r)}Vi(e,t)}},vs=function(){return ae},gs=function(e,t){var n=ae;try{return ae=e,t()}finally{ae=n}},ao=function(e,t,n){switch(t){case"input":if(hn(e,n),t=n.name,n.type==="radio"&&t!=null){for(n=e;n.parentNode;)n=n.parentNode;for(n=n.querySelectorAll("input[name="+JSON.stringify(""+t)+'][type="radio"]'),t=0;t<n.length;t++){var r=n[t];if(r!==e&&r.form===e.form){var l=yl(r);if(!l)throw Error(m(90));pn(r),hn(r,l)}}}break;case"textarea":Kr(e,n);break;case"select":t=n.value,t!=null&&Vt(e,!!n.multiple,t,!1)}},ts=Mi,ns=Pn;var Qd={usingClientEntryPoint:!1,Events:[Tr,Qn,yl,bi,es,Mi]},Qr={findFiberByHostInstance:wn,bundleType:0,version:"18.3.1",rendererPackageName:"react-dom"},Vd={bundleType:Qr.bundleType,version:Qr.version,rendererPackageName:Qr.rendererPackageName,rendererConfig:Qr.rendererConfig,overrideHookState:null,overrideHookStateDeletePath:null,overrideHookStateRenamePath:null,overrideProps:null,overridePropsDeletePath:null,overridePropsRenamePath:null,setErrorHandler:null,setSuspenseHandler:null,scheduleUpdate:null,currentDispatcherRef:Se.ReactCurrentDispatcher,findHostInstanceByFiber:function(e){return e=is(e),e===null?null:e.stateNode},findFiberByHostInstance:Qr.findFiberByHostInstance||Bd,findHostInstancesForRefresh:null,scheduleRefresh:null,scheduleRoot:null,setRefreshHandler:null,getCurrentFiber:null,reconcilerVersion:"18.3.1-next-f1338f8080-20240426"};if(typeof __REACT_DEVTOOLS_GLOBAL_HOOK__<"u"){var no=__REACT_DEVTOOLS_GLOBAL_HOOK__;if(!no.isDisabled&&no.supportsFiber)try{br=no.inject(Vd),Pt=no}catch{}}return et.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED=Qd,et.createPortal=function(e,t){var n=2<arguments.length&&arguments[2]!==void 0?arguments[2]:null;if(!Ki(t))throw Error(m(200));return Ud(e,t,null,n)},et.createRoot=function(e,t){if(!Ki(e))throw Error(m(299));var n=!1,r="",l=Iu;return t!=null&&(t.unstable_strictMode===!0&&(n=!0),t.identifierPrefix!==void 0&&(r=t.identifierPrefix),t.onRecoverableError!==void 0&&(l=t.onRecoverableError)),t=Qi(e,1,!1,null,null,n,!1,r,l),e[Dt]=t.current,Er(e.nodeType===8?e.parentNode:e),new Hi(t)},et.findDOMNode=function(e){if(e==null)return null;if(e.nodeType===1)return e;var t=e._reactInternals;if(t===void 0)throw typeof e.render=="function"?Error(m(188)):(e=Object.keys(e).join(","),Error(m(268,e)));return e=is(t),e=e===null?null:e.stateNode,e},et.flushSync=function(e){return Pn(e)},et.hydrate=function(e,t,n){if(!eo(t))throw Error(m(200));return to(null,e,t,!0,n)},et.hydrateRoot=function(e,t,n){if(!Ki(e))throw Error(m(405));var r=n!=null&&n.hydratedSources||null,l=!1,o="",i=Iu;if(n!=null&&(n.unstable_strictMode===!0&&(l=!0),n.identifierPrefix!==void 0&&(o=n.identifierPrefix),n.onRecoverableError!==void 0&&(i=n.onRecoverableError)),t=Ru(t,null,e,1,n??null,l,!1,o,i),e[Dt]=t.current,Er(e),r)for(e=0;e<r.length;e++)n=r[e],l=n._getVersion,l=l(n._source),t.mutableSourceEagerHydrationData==null?t.mutableSourceEagerHydrationData=[n,l]:t.mutableSourceEagerHydrationData.push(n,l);return new bl(t)},et.render=function(e,t,n){if(!eo(t))throw Error(m(200));return to(null,e,t,!1,n)},et.unmountComponentAtNode=function(e){if(!eo(e))throw Error(m(40));return e._reactRootContainer?(Pn(function(){to(null,null,e,!1,function(){e._reactRootContainer=null,e[Dt]=null})}),!0):!1},et.unstable_batchedUpdates=Mi,et.unstable_renderSubtreeIntoContainer=function(e,t,n,r){if(!eo(n))throw Error(m(200));if(e==null||e._reactInternals===void 0)throw Error(m(38));return to(e,t,n,!1,r)},et.version="18.3.1-next-f1338f8080-20240426",et}var Uu;function Jd(){if(Uu)return Gi.exports;Uu=1;function c(){if(!(typeof __REACT_DEVTOOLS_GLOBAL_HOOK__>"u"||typeof __REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE!="function"))try{__REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE(c)}catch(_){console.error(_)}}return c(),Gi.exports=Zd(),Gi.exports}var Bu;function bd(){if(Bu)return ro;Bu=1;var c=Jd();return ro.createRoot=c.createRoot,ro.hydrateRoot=c.hydrateRoot,ro}var ef=bd(),Xu=function(c,_,m,T){var E;_[0]=0;for(var x=1;x<_.length;x++){var M=_[x++],G=_[x]?(_[0]|=M?1:2,m[_[x++]]):_[++x];M===3?T[0]=G:M===4?T[1]=Object.assign(T[1]||{},G):M===5?(T[1]=T[1]||{})[_[++x]]=G:M===6?T[1][_[++x]]+=G+"":M?(E=c.apply(G,Xu(c,G,m,["",null])),T.push(E),G[0]?_[0]|=2:(_[x-2]=0,_[x]=E)):T.push(G)}return T},Wu=new Map;function tf(c){var _=Wu.get(this);return _||(_=new Map,Wu.set(this,_)),(_=Xu(this,_.get(c)||(_.set(c,_=(function(m){for(var T,E,x=1,M="",G="",C=[0],ue=function(Z){x===1&&(Z||(M=M.replace(/^\s*\n\s*|\s*\n\s*$/g,"")))?C.push(0,Z,M):x===3&&(Z||M)?(C.push(3,Z,M),x=2):x===2&&M==="..."&&Z?C.push(4,Z,0):x===2&&M&&!Z?C.push(5,0,!0,M):x>=5&&((M||!Z&&x===5)&&(C.push(x,0,M,E),x=6),Z&&(C.push(x,Z,0,E),x=6)),M=""},ye=0;ye<m.length;ye++){ye&&(x===1&&ue(),ue(ye));for(var le=0;le<m[ye].length;le++)T=m[ye][le],x===1?T==="<"?(ue(),C=[C],x=3):M+=T:x===4?M==="--"&&T===">"?(x=1,M=""):M=T+M[0]:G?T===G?G="":M+=T:T==='"'||T==="'"?G=T:T===">"?(ue(),x=1):x&&(T==="="?(x=5,E=M,M=""):T==="/"&&(x<5||m[ye][le+1]===">")?(ue(),x===3&&(C=C[0]),x=C,(C=C[0]).push(2,0,x),x=0):T===" "||T==="	"||T===`
`||T==="\r"?(ue(),x=2):M+=T),x===3&&M==="!--"&&(x=4,C=C[0])}return ue(),C})(c)),_),arguments,[])).length>1?_:_[0]}const nf="https://publicfiles-envolve.s3.us-west-2.amazonaws.com/logo_black.svg",rf="https://publicfiles-envolve.s3.us-west-2.amazonaws.com/logo_white.svg",lf=7e3,of=3e3,qu="whatsmeow-dashboard-theme",K=tf.bind(Gd.createElement);function Ee(c){return encodeURIComponent(c)}function sf(){const c=localStorage.getItem(qu);return c==="light"||c==="dark"?c:window.matchMedia&&window.matchMedia("(prefers-color-scheme: dark)").matches?"dark":"light"}async function Ce(c,_={}){const m={..._},T=new Headers(m.headers||{});m.body&&!(m.body instanceof FormData)&&!T.has("Content-Type")&&T.set("Content-Type","application/json");const E=await fetch(c,{...m,headers:T});if(!E.ok){let M=`${E.status} ${E.statusText}`;try{if((E.headers.get("content-type")||"").includes("application/json")){const C=await E.json();M=C?.error?.message||C?.message||C?.error||JSON.stringify(C)}else{const C=await E.text();C&&(M=C)}}catch{}throw new Error(M)}return(E.headers.get("content-type")||"").includes("application/json")?E.json():null}function af(){return`${Date.now()}-${Math.random().toString(16).slice(2)}`}function Zi(){return{reject_calls:!1,reject_calls_message:"",mark_read_on_message:!1,ignore_status_broadcast:!1,ignore_newsletters:!1,outgoing_messages_coex:!1,history_sync_enabled:!1,history_sync_days:"7",always_online:!1,always_online_interval_seconds:"60",proxy_enabled:!1,proxy_protocol:"",proxy_host:"",proxy_port:"",proxy_username:"",proxy_password:""}}function Qu(c){const _=c?.config||{};return{reject_calls:!!_.reject_calls,reject_calls_message:_.reject_calls_message||"",mark_read_on_message:!!_.mark_read_on_message,ignore_status_broadcast:!!_.ignore_status_broadcast,ignore_newsletters:!!_.ignore_newsletters,outgoing_messages_coex:!!_.outgoing_messages_coex,history_sync_enabled:!!_.history_sync_enabled,history_sync_days:String(Number.isFinite(Number(_.history_sync_days))&&Number(_.history_sync_days)>0?Number(_.history_sync_days):7),always_online:!!_.always_online,always_online_interval_seconds:String(Number.isFinite(Number(_.always_online_interval_seconds))?Number(_.always_online_interval_seconds):60),proxy_enabled:!!_.proxy_enabled,proxy_protocol:_.proxy_protocol||"",proxy_host:_.proxy_host||"",proxy_port:_.proxy_port||"",proxy_username:_.proxy_username||"",proxy_password:_.proxy_password||""}}function Vu(c){const _=parseInt(c.timeout_ms,10);return{id:c.id.trim(),name:c.name.trim(),url:c.url.trim(),url_absolute:c.url_absolute.trim(),header:c.header.trim(),token:c.token.trim(),timeout_ms:Number.isFinite(_)&&_>0?_:5e3,enabled:c.enabled,append_session:c.append_session,send_incoming_messages:c.send_incoming_messages,send_outgoing_messages:c.send_outgoing_messages,send_update_messages:c.send_update_messages,send_group_messages:c.send_group_messages,send_newsletter_messages:c.send_newsletter_messages}}function uf(c,_){const m=parseInt(_,10);return{history_sync_enabled:!!c,history_sync_days:Number.isFinite(m)&&m>0?m:0}}function dn(){return{id:"",name:"",url:"",url_absolute:"",header:"",token:"",timeout_ms:"5000",enabled:!0,append_session:!0,send_incoming_messages:!0,send_outgoing_messages:!0,send_update_messages:!0,send_group_messages:!0,send_newsletter_messages:!1}}function cf(c){return{id:c?.id||"",name:c?.name||"",url:c?.url||"",url_absolute:c?.url_absolute||"",header:c?.header||"",token:c?.token||"",timeout_ms:String(Number.isFinite(Number(c?.timeout_ms))?Number(c.timeout_ms):5e3),enabled:!!c?.enabled,append_session:!!c?.append_session,send_incoming_messages:!!c?.send_incoming_messages,send_outgoing_messages:!!c?.send_outgoing_messages,send_update_messages:!!c?.send_update_messages,send_group_messages:!!c?.send_group_messages,send_newsletter_messages:!!c?.send_newsletter_messages}}function lr(){return{id:"",name:"",language:"pt_BR",category:"",header_type:"",header_text:"",header_link:"",header_id:"",header_filename:"",body:"",footer:""}}function df(c){const _=c?.header||{};return{id:c?.id||"",name:c?.name||"",language:c?.language||"pt_BR",category:c?.category||"",header_type:_?.type||"",header_text:_?.text||"",header_link:_?.link||"",header_id:_?.id||"",header_filename:_?.filename||"",body:c?.body||"",footer:c?.footer||""}}function ff(c){return{id:c.id.trim(),name:c.name.trim(),language:c.language.trim(),category:c.category.trim(),header:{type:c.header_type.trim(),text:c.header_text.trim(),link:c.header_link.trim(),id:c.header_id.trim(),filename:c.header_filename.trim()},body:c.body.trim(),footer:c.footer.trim()}}function fn(c){const _=String(c?.status||"").trim().toLowerCase();return _==="connected"?"online":_||(c?.connected?"online":"offline")}function Zu(c){return c==="online"?"online":c==="connecting"?"connecting":"offline"}function Hu(c,_){const m=fn(c),T=fn(_),E=(c?.assigned?0:4)+(m==="online"?0:m==="connecting"?1:2),x=(_?.assigned?0:4)+(T==="online"?0:T==="connecting"?1:2);return E!==x?E-x:String(c?.session_id||"").localeCompare(String(_?.session_id||""),"pt-BR")}function pf(c){return c?new Intl.DateTimeFormat("pt-BR",{hour:"2-digit",minute:"2-digit",second:"2-digit"}).format(c):"Nunca"}function mf(c,_=56){return c?c.length>_?`${c.slice(0,_)}…`:c:"—"}function Ku(c){return String(c||"").replace(/\D+/g,"")}function hf(c){return c.key==="Enter"||c.key===" "}function lo({label:c,value:_,caption:m,tone:T="neutral"}){return K`
    <article className=${`stat-card stat-card--${T}`}>
      <p className="stat-card__label">${c}</p>
      <strong className="stat-card__value">${_}</strong>
      <span className="stat-card__caption">${m}</span>
    </article>
  `}function re({label:c,value:_,onInput:m,placeholder:T="",type:E="text",min:x,step:M}){return K`
    <label className="field">
      <span>${c}</span>
      <input
        type=${E}
        value=${_}
        placeholder=${T}
        min=${x}
        step=${M}
        onInput=${G=>m(G.target.value)}
      />
    </label>
  `}function oo({label:c,value:_,onInput:m,placeholder:T="",rows:E=4}){return K`
    <label className="field">
      <span>${c}</span>
      <textarea
        rows=${E}
        value=${_}
        placeholder=${T}
        onInput=${x=>m(x.target.value)}
      ></textarea>
    </label>
  `}function De({label:c,description:_="",checked:m,onChange:T}){return K`
    <label className="check-field">
      <input
        type="checkbox"
        checked=${m}
        onChange=${E=>T(E.target.checked)}
      />
      <span>
        <strong>${c}</strong>
        ${_?K`<small>${_}</small>`:null}
      </span>
    </label>
  `}function vf({active:c,label:_,count:m,onClick:T}){return K`
    <button
      type="button"
      className=${`filter-chip ${c?"filter-chip--active":""}`}
      onClick=${T}
    >
      <span>${_}</span>
      <strong>${m}</strong>
    </button>
  `}function Wt({label:c,tone:_="ghost",onClick:m,disabled:T=!1}){return K`
    <button
      type="button"
      className=${`mini-button mini-button--${_}`}
      disabled=${T}
      onClick=${E=>{E.stopPropagation(),m()}}
    >
      ${c}
    </button>
  `}function gf({session:c,selected:_,busyMap:m,onSelect:T,onSessionAction:E,onOpenTab:x,onShowQR:M}){const G=fn(c),C=Zu(G),ue=c.assigned?"Desatribuir":"Atribuir",ye=c.assigned?"unassign":"assign";return K`
    <article
      className=${`session-card ${_?"session-card--selected":""}`}
      role="button"
      tabIndex="0"
      onClick=${()=>T(c.session_id)}
      onKeyDown=${le=>{hf(le)&&(le.preventDefault(),T(c.session_id))}}
    >
      <div className="session-card__top">
        <div>
          <p className="session-card__eyebrow">Instância</p>
          <h3 className="session-card__title">${c.session_id}</h3>
        </div>
        <span className=${`status-pill status-pill--${C}`}>${G}</span>
      </div>

      <div className="tag-row">
        <span className="tag">${c.connected?"Conectada":"Desconectada"}</span>
        <span className="tag">${c.assigned?"Atribuída":"Sem atribuição"}</span>
        <span className="tag">${c.has_qr?"QR disponível":"Sem QR"}</span>
      </div>

      <div className="session-card__meta">
        <span>Login: ${c.logged_in?"ativo":"pendente"}</span>
        <span>Device: ${c.device_jid?mf(c.device_jid,30):"—"}</span>
        <span>WA Web: ${c.whatsapp_web_version||"—"}</span>
      </div>

      <div className="session-card__actions">
        <${Wt}
          label="QR"
          onClick=${()=>M(c.session_id)}
          disabled=${!c.has_qr}
        />
        <${Wt}
          label="Conectar"
          onClick=${()=>E("connect",c.session_id)}
          disabled=${!!m[`connect:${c.session_id}`]}
        />
        <${Wt}
          label="Recarregar"
          onClick=${()=>E("reload",c.session_id)}
          disabled=${!!m[`reload:${c.session_id}`]}
        />
        <${Wt}
          label="Desconectar"
          tone="danger"
          onClick=${()=>E("disconnect",c.session_id)}
          disabled=${!!m[`disconnect:${c.session_id}`]}
        />
        <${Wt}
          label="Excluir"
          tone="danger"
          onClick=${()=>E("delete",c.session_id)}
          disabled=${!!m[`delete:${c.session_id}`]}
        />
        <${Wt}
          label=${ue}
          onClick=${()=>E(ye,c.session_id)}
          disabled=${!!m[`${ye}:${c.session_id}`]}
        />
        <${Wt}
          label="Config"
          onClick=${()=>x(c.session_id,"config")}
        />
        <${Wt}
          label="Hooks"
          onClick=${()=>x(c.session_id,"webhooks")}
        />
        <${Wt}
          label="Templates"
          onClick=${()=>x(c.session_id,"templates")}
        />
      </div>
    </article>
  `}function at({label:c,value:_,code:m=!1}){return K`
    <div className="detail-item">
      <span>${c}</span>
      ${m?K`<code>${_}</code>`:K`<strong>${_}</strong>`}
    </div>
  `}function yf({sessionId:c,onClose:_}){const[m,T]=ie.useState(Date.now());return ie.useEffect(()=>{if(!c)return;T(Date.now());const E=window.setInterval(()=>T(Date.now()),of),x=M=>{M.key==="Escape"&&_()};return window.addEventListener("keydown",x),()=>{window.clearInterval(E),window.removeEventListener("keydown",x)}},[c]),c?K`
    <div
      className="modal-backdrop"
      onClick=${E=>{E.target===E.currentTarget&&_()}}
    >
      <div className="modal-card">
        <div className="modal-card__header">
          <div>
            <p className="eyebrow">Pareamento</p>
            <h3>QR da sessão ${c}</h3>
          </div>
          <button type="button" className="ghost-button" onClick=${_}>Fechar</button>
        </div>
        <img
          className="modal-card__image"
          src=${`/sessions/${Ee(c)}/qr?ts=${m}`}
          alt=${`QR da sessão ${c}`}
        />
        <p className="hint-text">
          O QR atualiza automaticamente a cada poucos segundos enquanto a sessão estiver em pareamento.
        </p>
      </div>
    </div>
  `:null}function _f({item:c,onEdit:_,onDelete:m}){return K`
    <article className="webhook-card">
      <div className="webhook-card__top">
        <div>
          <strong>${c.id}</strong>
          <p>${c.name||"Webhook sem nome"}</p>
        </div>
        <div className="inline-actions">
          <button type="button" className="text-button" onClick=${_}>Editar</button>
          <button type="button" className="text-button text-button--danger" onClick=${m}>
            Excluir
          </button>
        </div>
      </div>

      <code className="webhook-card__url">${c.url_absolute||c.url||"—"}</code>

      <div className="tag-row">
        <span className="tag">${c.enabled?"Ativo":"Desligado"}</span>
        <span className="tag">${c.append_session?"Anexa sessão":"URL fixa"}</span>
        <span className="tag">${c.send_incoming_messages?"Incoming":"Sem incoming"}</span>
        <span className="tag">${c.send_outgoing_messages?"Outgoing":"Sem outgoing"}</span>
        <span className="tag">${c.send_update_messages?"Update":"Sem update"}</span>
      </div>
    </article>
  `}function wf({item:c,onEdit:_,onDelete:m}){const T=c?.header||{},E=T.type||(T.text?"text":"sem cabeçalho");return K`
    <article className="webhook-card">
      <div className="webhook-card__top">
        <div>
          <strong>${c.name}</strong>
          <p>${c.language||"idioma padrão"}${c.category?` · ${c.category}`:""}</p>
        </div>
        <div className="inline-actions">
          <button type="button" className="text-button" onClick=${_}>Editar</button>
          <button type="button" className="text-button text-button--danger" onClick=${m}>
            Excluir
          </button>
        </div>
      </div>

      <code className="webhook-card__url">${c.body||"—"}</code>

      <div className="tag-row">
        <span className="tag">Header: ${E}</span>
        <span className="tag">${c.footer?"Com footer":"Sem footer"}</span>
        <span className="tag">${c.id}</span>
      </div>
    </article>
  `}function kf(){const[c,_]=ie.useState(sf),[m,T]=ie.useState(null),[E,x]=ie.useState([]),[M,G]=ie.useState("home"),[C,ue]=ie.useState(""),[ye,le]=ie.useState("quick"),[Z,ut]=ie.useState(!0),[Oe,fe]=ie.useState(!0),[oe,vt]=ie.useState(null),[tt,nt]=ie.useState(""),[Se,rt]=ie.useState("all"),[Fe,Me]=ie.useState([]),[Be,xt]=ie.useState({}),[gt,yt]=ie.useState(""),[J,me]=ie.useState({id:"",assign:!0,connect:!0,history_sync_enabled:!1,history_sync_days:"7",webhook:dn()}),[Pe,We]=ie.useState({session:"",to:"",body:""}),[j,Q]=ie.useState({loading:!1,saving:!1,snapshot:null,form:Zi()}),[P,U]=ie.useState({phone:"",code:"",loading:!1}),[$,f]=ie.useState({loading:!1,saving:!1,items:[],form:dn()}),[y,D]=ie.useState({loading:!1,saving:!1,items:[],form:lr()}),Y=ie.useRef(null);function z(a,u="info"){const L=af();Me(V=>[...V,{id:L,message:a,tone:u}]),window.setTimeout(()=>{Me(V=>V.filter(de=>de.id!==L))},3400)}function X(a,u){xt(L=>{const V={...L};return u?V[a]=!0:delete V[a],V})}async function b({silent:a=!1}={}){a||ut(!0);try{const[u,L]=await Promise.all([Ce("/api/v1/worker"),Ce("/api/v1/sessions")]);T(u),x(Array.isArray(L?.sessions)?L.sessions:[]),vt(new Date)}catch(u){z(u.message,"error")}finally{a||ut(!1)}}Y.current=b;async function ee(a){const u=String(a||"").trim();if(u){Q(L=>({...L,loading:!0}));try{const L=await Ce(`/api/v1/sessions/${Ee(u)}/config`);Q(V=>({...V,loading:!1,snapshot:L,form:Qu(L)}))}catch(L){Q(V=>({...V,loading:!1})),z(L.message,"error")}}}async function ce(a,{resetForm:u=!1}={}){const L=String(a||"").trim();if(L){f(V=>({...V,loading:!0}));try{const V=await Ce(`/api/v1/sessions/${Ee(L)}/webhooks`);f(de=>({...de,loading:!1,items:Array.isArray(V?.webhooks)?V.webhooks:[],form:u?dn():de.form}))}catch(V){f(de=>({...de,loading:!1})),z(V.message,"error")}}}async function je(a,{resetForm:u=!1}={}){const L=String(a||"").trim();if(L){D(V=>({...V,loading:!0}));try{const V=await Ce(`/api/v1/sessions/${Ee(L)}/templates`);D(de=>({...de,loading:!1,items:Array.isArray(V?.templates)?V.templates:[],form:u?lr():de.form}))}catch(V){D(de=>({...de,loading:!1})),z(V.message,"error")}}}ie.useEffect(()=>{document.documentElement.dataset.theme=c,localStorage.setItem(qu,c)},[c]),ie.useEffect(()=>{document.title=m?.ui_title||"WhatsMeow Control Center"},[m]),ie.useEffect(()=>{b()},[]),ie.useEffect(()=>{if(!Oe)return;const a=window.setInterval(()=>{!document.hidden&&Y.current&&Y.current({silent:!0})},lf);return()=>window.clearInterval(a)},[Oe]),ie.useEffect(()=>{if(E.length===0){ue("");return}if(C&&E.some(u=>u.session_id===C))return;const a=[...E].sort(Hu)[0];a&&ue(a.session_id)},[E,C]),ie.useEffect(()=>{if(!C){Q({loading:!1,saving:!1,snapshot:null,form:Zi()}),U({phone:"",code:"",loading:!1}),f({loading:!1,saving:!1,items:[],form:dn()}),D({loading:!1,saving:!1,items:[],form:lr()});return}We(u=>({...u,session:C})),U({phone:Ku(C),code:"",loading:!1});let a=!0;return Q({loading:!0,saving:!1,snapshot:null,form:Zi()}),f({loading:!0,saving:!1,items:[],form:dn()}),D({loading:!0,saving:!1,items:[],form:lr()}),Promise.all([Ce(`/api/v1/sessions/${Ee(C)}/config`),Ce(`/api/v1/sessions/${Ee(C)}/webhooks`),Ce(`/api/v1/sessions/${Ee(C)}/templates`)]).then(([u,L,V])=>{a&&(Q({loading:!1,saving:!1,snapshot:u,form:Qu(u)}),f({loading:!1,saving:!1,items:Array.isArray(L?.webhooks)?L.webhooks:[],form:dn()}),D({loading:!1,saving:!1,items:Array.isArray(V?.templates)?V.templates:[],form:lr()}))}).catch(u=>{a&&(Q(L=>({...L,loading:!1})),f(L=>({...L,loading:!1})),D(L=>({...L,loading:!1})),z(u.message,"error"))}),()=>{a=!1}},[C]),m?.ui_title;const se=E.find(a=>a.session_id===C)||null,pn=E.filter(a=>a.connected||fn(a)==="online").length,mn=E.filter(a=>a.has_qr).length,Qt=E.filter(a=>a.assigned).length,Vr=[...E].sort(Hu),or=tt.trim().toLowerCase(),hn=Vr.filter(a=>Se==="assigned"&&!a.assigned||Se==="online"&&!(a.connected||fn(a)==="online")||Se==="qr"&&!a.has_qr?!1:or?[a.session_id,a.device_jid,fn(a),a.assigned?"atribuida":"livre"].filter(Boolean).join(" ").toLowerCase().includes(or):!0),Ln=se?fn(se):"",ir=!!(se&&se.has_qr&&Ln==="connecting");async function vn(a){a.preventDefault();const u=J.id.trim();if(!u){z("Informe o ID da sessão.","error");return}const L=parseInt(J.history_sync_days,10);if(J.history_sync_enabled){if(!Number.isFinite(L)||L<=0){z("Informe a quantidade de dias do histórico.","error");return}if(!J.webhook.url.trim()&&!J.webhook.url_absolute.trim()){z("Informe o webhook antes de ativar a sincronização de histórico.","error");return}}const V=`create:${u}`;X(V,!0);try{const de=uf(J.history_sync_enabled,J.history_sync_days),_t={id:u,assign:J.assign,connect:J.connect,config:de};J.history_sync_enabled&&(_t.webhook=Vu(J.webhook)),await Ce("/api/v1/sessions",{method:"POST",body:JSON.stringify(_t)}),z(`Sessão ${u} solicitada com sucesso.`,"success"),me({id:"",assign:!0,connect:!0,history_sync_enabled:!1,history_sync_days:"7",webhook:dn()}),ue(u),le("overview"),await b({silent:!0})}catch(de){z(de.message,"error")}finally{X(V,!1)}}async function Vt(a){a.preventDefault();const u=Pe.session.trim(),L=Pe.to.trim(),V=Pe.body.trim();if(!u||!L||!V){z("Preencha sessão, destino e mensagem.","error");return}const de=`send:${u}`;X(de,!0);try{await Ce(`/api/v1/sessions/${Ee(u)}/messages`,{method:"POST",body:JSON.stringify({type:"text",to:L,text:{body:V}})}),z(`Mensagem enviada pela sessão ${u}.`,"success"),We(_t=>({..._t,body:""}))}catch(_t){z(_t.message,"error")}finally{X(de,!1)}}async function It(a,u){const L=String(u||"").trim();if(!L)return;if(a==="disconnect"||a==="delete"){const de=a==="delete"?`Excluir permanentemente a sessão "${L}"? Isso remove o diretório da sessão, a atribuição, overrides e webhooks.`:`Desconectar a sessão "${L}"?`;if(!window.confirm(de))return}const V=`${a}:${L}`;X(V,!0);try{a==="assign"?await Ce("/api/v1/worker/assignments",{method:"POST",body:JSON.stringify({session_id:L})}):a==="unassign"?await Ce(`/api/v1/worker/assignments/${Ee(L)}`,{method:"DELETE"}):a==="delete"?await Ce(`/api/v1/sessions/${Ee(L)}`,{method:"DELETE"}):a==="qr"?yt(L):await Ce(`/sessions/${Ee(L)}/${a}`,{method:"POST"}),a!=="qr"&&(z(a==="delete"?`Sessão ${L} excluída.`:`Ação "${a}" aplicada em ${L}.`,"success"),await b({silent:!0}))}catch(de){z(de.message,"error")}finally{X(V,!1)}}function Hr(a,u){ue(a),G("instances"),le(u)}async function Kr(a){if(a.preventDefault(),!C)return;const u=parseInt(j.form.always_online_interval_seconds,10);if(!Number.isFinite(u)||u<=0){z("O intervalo de always-online deve ser maior que zero.","error");return}const L=parseInt(j.form.history_sync_days,10);if(j.form.history_sync_enabled&&(!Number.isFinite(L)||L<=0)){z("A quantidade de dias do histórico deve ser maior que zero.","error");return}if(j.form.proxy_enabled){if(!j.form.proxy_host.trim()){z("Informe o host do proxy.","error");return}if(!j.form.proxy_port.trim()){z("Informe a porta do proxy.","error");return}}Q(V=>({...V,saving:!0}));try{await Ce(`/api/v1/sessions/${Ee(C)}/config`,{method:"PUT",body:JSON.stringify({reject_calls:j.form.reject_calls,reject_calls_message:j.form.reject_calls_message.trim(),mark_read_on_message:j.form.mark_read_on_message,ignore_status_broadcast:j.form.ignore_status_broadcast,ignore_newsletters:j.form.ignore_newsletters,outgoing_messages_coex:j.form.outgoing_messages_coex,history_sync_enabled:j.form.history_sync_enabled,history_sync_days:j.form.history_sync_enabled?L:0,always_online:j.form.always_online,always_online_interval_seconds:u,proxy_enabled:j.form.proxy_enabled,proxy_protocol:j.form.proxy_protocol.trim(),proxy_host:j.form.proxy_host.trim(),proxy_port:j.form.proxy_port.trim(),proxy_username:j.form.proxy_username.trim(),proxy_password:j.form.proxy_password})}),z(`Configuração da sessão ${C} salva.`,"success"),await ee(C)}catch(V){z(V.message,"error")}finally{Q(V=>({...V,saving:!1}))}}async function Yr(a){if(a.preventDefault(),!C)return;const u=Ku(P.phone);if(!u){z("Informe o telefone para gerar o código de pareamento.","error");return}U(L=>({...L,phone:u,code:"",loading:!0}));try{const L=await Ce(`/api/v1/sessions/${Ee(C)}/pairing-code`,{method:"POST",body:JSON.stringify({phone:u})}),V=String(L?.pairing_code||"").trim();if(!V)throw new Error("pairing code not returned by server");U(de=>({...de,phone:u,code:V,loading:!1})),z(`Código de pareamento gerado para ${C}.`,"success")}catch(L){U(V=>({...V,loading:!1})),z(L.message,"error")}}async function Gr(){if(C&&window.confirm(`Remover o override da sessão "${C}"?`)){Q(a=>({...a,saving:!0}));try{await Ce(`/api/v1/sessions/${Ee(C)}/config`,{method:"DELETE"}),z(`Override da sessão ${C} removido.`,"success"),await ee(C)}catch(a){z(a.message,"error")}finally{Q(a=>({...a,saving:!1}))}}}async function sr(a){if(a.preventDefault(),!C)return;const u=parseInt($.form.timeout_ms,10);if(!$.form.url.trim()&&!$.form.url_absolute.trim()){z("Informe a URL base ou a URL absoluta do webhook.","error");return}if(!Number.isFinite(u)||u<=0){z("O timeout do webhook deve ser maior que zero.","error");return}f(L=>({...L,saving:!0}));try{await Ce(`/api/v1/sessions/${Ee(C)}/webhooks`,{method:"POST",body:JSON.stringify(Vu($.form))}),z(`Webhook salvo para a sessão ${C}.`,"success"),await ce(C,{resetForm:!0})}catch(L){z(L.message,"error")}finally{f(L=>({...L,saving:!1}))}}async function In(a){if(!(!C||!a)&&window.confirm(`Excluir o webhook "${a}" da sessão "${C}"?`)){f(u=>({...u,saving:!0}));try{await Ce(`/api/v1/sessions/${Ee(C)}/webhooks/${Ee(a)}`,{method:"DELETE"}),z(`Webhook ${a} excluído.`,"success"),await ce(C,{resetForm:!0})}catch(u){z(u.message,"error")}finally{f(u=>({...u,saving:!1}))}}}async function Xr(a){if(a.preventDefault(),!!C){if(!y.form.name.trim()){z("Informe o nome do template.","error");return}if(!y.form.body.trim()){z("Informe o corpo do template.","error");return}D(u=>({...u,saving:!0}));try{await Ce(`/api/v1/sessions/${Ee(C)}/templates`,{method:"POST",body:JSON.stringify(ff(y.form))}),z(`Template salvo para a sessão ${C}.`,"success"),await je(C,{resetForm:!0})}catch(u){z(u.message,"error")}finally{D(u=>({...u,saving:!1}))}}}async function gn(a){if(!(!C||!a)&&window.confirm(`Excluir o template "${a}" da sessão "${C}"?`)){D(u=>({...u,saving:!0}));try{await Ce(`/api/v1/sessions/${Ee(C)}/templates/${Ee(a)}`,{method:"DELETE"}),z(`Template ${a} excluído.`,"success"),await je(C,{resetForm:!0})}catch(u){z(u.message,"error")}finally{D(u=>({...u,saving:!1}))}}}const yn=[{id:"all",label:"Todas",count:E.length},{id:"assigned",label:"Atribuídas",count:Qt},{id:"online",label:"Online",count:pn},{id:"qr",label:"Com QR",count:mn}];return K`
    <div className="page">
      <main className="dashboard-shell">
        <aside className="system-sidebar">
          <div className="system-sidebar__brand">
            <img
              className="system-sidebar__logo"
              src=${c==="dark"?rf:nf}
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
            <span>${pn} online</span>
            <span>${Qt} atribuída(s)</span>
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
              ${oe?`Atualizado às ${pf(oe)}`:"Carregando dados"}
            </span>
            <button
              type="button"
              className=${`ghost-button ${Oe?"ghost-button--active":""}`}
              onClick=${()=>fe(a=>!a)}
            >
              ${Oe?"Auto ligado":"Auto desligado"}
            </button>
            <button
              type="button"
              className="solid-button"
              disabled=${Z}
              onClick=${()=>b()}
            >
              ${Z?"Atualizando...":"Atualizar"}
            </button>
            <button
              type="button"
              className="ghost-button"
              onClick=${()=>_(a=>a==="dark"?"light":"dark")}
            >
              ${c==="dark"?"Tema claro":"Tema escuro"}
            </button>
          </div>
        </section>

        ${M==="home"?K`
        <section className="summary-grid">
          <${lo}
            label="Worker ativo"
            value=${m?.worker_id||"—"}
            caption=${m?.worker_strict?"Modo strict habilitado":"Aceita sessões fora da lista"}
            tone="brand"
          />
          <${lo}
            label="Sessões atribuídas"
            value=${String(Qt)}
            caption="visíveis para este worker"
            tone="accent"
          />
          <${lo}
            label="Sessões online"
            value=${String(pn)}
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
              <${at} label="Total" value=${String(E.length)} />
              <${at} label="Online" value=${String(pn)} />
              <${at} label="Com QR" value=${String(mn)} />
              <${at} label="Atribuídas" value=${String(Qt)} />
            </div>
            <div className="inline-actions inline-actions--wrap">
              <button type="button" className="solid-button" onClick=${()=>G("instances")}>
                Abrir instâncias
              </button>
              <button type="button" className="ghost-button" onClick=${()=>b()}>
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
              ${(m?.amqp_bindings||[]).length?m.amqp_bindings.map(a=>K`<span key=${a} className="tag">${a}</span>`):K`<span className="empty-inline">Nenhum binding configurado.</span>`}
            </div>
          </section>
        </section>
          `:null}

        ${M==="instances"?K`
        <section className="workspace-grid">
          <div className="workspace-column workspace-column--sessions">
            <section className="panel">
              <div className="panel__header panel__header--space-between">
                <div>
                  <p className="eyebrow">Instâncias</p>
                  <h2>Instâncias</h2>
                </div>
                <span className="panel__hint">
                  ${hn.length} de ${E.length}
                </span>
              </div>

              <details className="create-instance">
                <summary>Nova instância</summary>
                <form className="stack-form" onSubmit=${vn}>
                  <${re}
                    label="ID da sessão"
                    value=${J.id}
                    placeholder="Ex: 5562999999999"
                    onInput=${a=>me(u=>({...u,id:a}))}
                  />
                  <div className="check-grid check-grid--single">
                    <${De}
                      label="Atribuir ao worker atual"
                      checked=${J.assign}
                      onChange=${a=>me(u=>({...u,assign:a}))}
                    />
                    <${De}
                      label="Conectar imediatamente"
                      checked=${J.connect}
                      onChange=${a=>me(u=>({...u,connect:a}))}
                    />
                    <${De}
                      label="Sincronizar histórico"
                      description="Configure o webhook antes da conexão para receber o HistorySync."
                      checked=${J.history_sync_enabled}
                      onChange=${a=>me(u=>({...u,history_sync_enabled:a}))}
                    />
                  </div>

                  ${J.history_sync_enabled?K`
                        <div className="form-grid">
                          <${re}
                            label="Dias de histórico"
                            value=${J.history_sync_days}
                            type="number"
                            min="1"
                            step="1"
                            onInput=${a=>me(u=>({...u,history_sync_days:a}))}
                          />
                          <${re}
                            label="Webhook ID"
                            value=${J.webhook.id}
                            placeholder="vazio = auto"
                            onInput=${a=>me(u=>({...u,webhook:{...u.webhook,id:a}}))}
                          />
                        </div>

                        <${re}
                          label="URL base do webhook"
                          value=${J.webhook.url}
                          placeholder="https://seu-app.com/webhooks/whatsapp"
                          onInput=${a=>me(u=>({...u,webhook:{...u.webhook,url:a}}))}
                        />
                        <${re}
                          label="URL absoluta do webhook"
                          value=${J.webhook.url_absolute}
                          placeholder="https://seu-app.com/webhooks/whatsapp/5562..."
                          onInput=${a=>me(u=>({...u,webhook:{...u.webhook,url_absolute:a}}))}
                        />

                        <div className="form-grid">
                          <${re}
                            label="Header auth"
                            value=${J.webhook.header}
                            placeholder="Authorization"
                            onInput=${a=>me(u=>({...u,webhook:{...u.webhook,header:a}}))}
                          />
                          <${re}
                            label="Token"
                            value=${J.webhook.token}
                            placeholder="Bearer xxxxx"
                            onInput=${a=>me(u=>({...u,webhook:{...u.webhook,token:a}}))}
                          />
                        </div>
                      `:null}
                  <button
                    type="submit"
                    className="solid-button solid-button--wide"
                    disabled=${!!Be[`create:${J.id.trim()}`]}
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
                  ${yn.map(a=>K`
                      <${vf}
                        key=${a.id}
                        active=${Se===a.id}
                        label=${a.label}
                        count=${a.count}
                        onClick=${()=>rt(a.id)}
                      />
                    `)}
                </div>
              </div>

              ${Z?K`
                    <div className="empty-state">
                      <strong>Carregando sessões…</strong>
                      <p>Consultando o worker e o catálogo de sessões disponíveis.</p>
                    </div>
                  `:hn.length===0?K`
                      <div className="empty-state">
                        <strong>Nenhuma sessão encontrada</strong>
                        <p>Ajuste o filtro ou crie uma nova instância.</p>
                      </div>
                    `:K`
                      <div className="session-grid">
                        ${hn.map(a=>K`
                            <${gf}
                              key=${a.session_id}
                              session=${a}
                              selected=${C===a.session_id}
                              busyMap=${Be}
                              onSelect=${u=>{ue(u),G("instances")}}
                              onSessionAction=${It}
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
              ${se?K`
                    <div className="panel__header">
                      <div>
                        <p className="eyebrow">Instância</p>
                        <h2>${se.session_id}</h2>
                      </div>
                      <div className="session-header-actions">
                        <span className=${`status-pill status-pill--${Zu(Ln)}`}>
                          ${Ln}
                        </span>
                        <button
                          type="button"
                          className="ghost-button"
                          onClick=${()=>It("connect",se.session_id)}
                        >
                          Conectar
                        </button>
                        <button
                          type="button"
                          className="ghost-button"
                          onClick=${()=>It("reload",se.session_id)}
                        >
                          Recarregar
                        </button>
                        <button
                          type="button"
                          className="ghost-button ghost-button--danger"
                          onClick=${()=>It("disconnect",se.session_id)}
                        >
                          Desconectar
                        </button>
                        <button
                          type="button"
                          className="ghost-button ghost-button--danger"
                          onClick=${()=>It("delete",se.session_id)}
                        >
                          Excluir
                        </button>
                      </div>
                    </div>

                    <div className="session-status-strip">
                      <${at} label="Worker" value=${se.assigned?m?.worker_id||"atribuída":"livre"} />
                      <${at} label="Conexão" value=${se.connected?"conectada":"desconectada"} />
                      <${at} label="Login" value=${se.logged_in?"ativo":"pendente"} />
                      <${at} label="WA Web" value=${se.whatsapp_web_version||"—"} code=${!!se.whatsapp_web_version} />
                    </div>

                    ${ir?K`
                          <div className="inline-qr">
                            <img
                              src=${`/sessions/${Ee(se.session_id)}/qr?ts=${Date.now()}`}
                              alt=${`QR da sessão ${se.session_id}`}
                            />
                            <div>
                              <strong>QR disponível</strong>
                              <span>Escaneie para concluir a conexão desta instância.</span>
                            </div>
                          </div>
                        `:null}

                    <div className="tab-row">
                      ${["quick","config","webhooks","templates"].map(a=>K`
                          <button
                            key=${a}
                            type="button"
                            className=${`tab-button ${ye===a?"tab-button--active":""}`}
                            onClick=${()=>le(a)}
                          >
                            ${a==="quick"?"Teste rápido":a==="config"?"Config":a==="webhooks"?"Webhook":"Templates"}
                          </button>
                        `)}
                    </div>

                    ${ye==="quick"?K`
                          <div className="detail-grid">
                            <${at}
                              label="Device JID"
                              value=${se.device_jid||"—"}
                              code=${!!se.device_jid}
                            />
                            <${at}
                              label="WA Web"
                              value=${se.whatsapp_web_version||"—"}
                              code=${!!se.whatsapp_web_version}
                            />
                            <${at}
                              label="Status"
                              value=${fn(se)}
                            />
                            <${at}
                              label="Worker"
                              value=${se.assigned?m?.worker_id||"atribuída":"livre"}
                            />
                            <${at}
                              label="Pairing"
                              value=${se.has_qr?"QR disponível":"sem QR"}
                            />
                          </div>

                          <form className="stack-form" onSubmit=${Vt}>
                            <${re}
                              label="Sessão"
                              value=${Pe.session}
                              placeholder="Selecione uma sessão ou informe manualmente"
                              onInput=${a=>We(u=>({...u,session:a}))}
                            />
                            <${re}
                              label="Destino"
                              value=${Pe.to}
                              placeholder="5562999999999"
                              onInput=${a=>We(u=>({...u,to:a}))}
                            />
                            <${oo}
                              label="Mensagem"
                              value=${Pe.body}
                              rows="4"
                              placeholder="Texto da mensagem"
                              onInput=${a=>We(u=>({...u,body:a}))}
                            />
                            <button
                              type="submit"
                              className="solid-button"
                              disabled=${!!Be[`send:${Pe.session.trim()}`]}
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

                            <${re}
                              label="Telefone para o pareamento"
                              value=${P.phone}
                              placeholder="5562999999999"
                              onInput=${a=>U(u=>({...u,phone:a,code:""}))}
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

                            ${P.code?K`
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

                    ${ye==="config"?K`
                          <form className="stack-form" onSubmit=${Kr}>
                            <div
                              className=${`callout ${j.snapshot?.has_override?"callout--accent":"callout--muted"}`}
                            >
                              <strong>
                                ${j.snapshot?.has_override?"Existe um override ativo para esta sessão.":"A sessão está usando os defaults do worker."}
                              </strong>
                              <p>
                                ${j.loading?"Carregando configuração efetiva…":"Edite os parâmetros operacionais e salve no Redis/arquivo da sessão."}
                              </p>
                            </div>

                            <div className="check-grid">
                              <${De}
                                label="Rejeitar ligações"
                                checked=${j.form.reject_calls}
                                onChange=${a=>Q(u=>({...u,form:{...u.form,reject_calls:a}}))}
                              />
                              <${De}
                                label="Marcar conversa como lida"
                                checked=${j.form.mark_read_on_message}
                                onChange=${a=>Q(u=>({...u,form:{...u.form,mark_read_on_message:a}}))}
                              />
                              <${De}
                                label="Ignorar status broadcast"
                                checked=${j.form.ignore_status_broadcast}
                                onChange=${a=>Q(u=>({...u,form:{...u.form,ignore_status_broadcast:a}}))}
                              />
                              <${De}
                                label="Ignorar newsletters"
                                checked=${j.form.ignore_newsletters}
                                onChange=${a=>Q(u=>({...u,form:{...u.form,ignore_newsletters:a}}))}
                              />
                              <${De}
                                label="Outgoing em smb_message_echoes"
                                checked=${j.form.outgoing_messages_coex}
                                onChange=${a=>Q(u=>({...u,form:{...u.form,outgoing_messages_coex:a}}))}
                              />
                              <${De}
                                label="Manter sempre online"
                                checked=${j.form.always_online}
                                onChange=${a=>Q(u=>({...u,form:{...u.form,always_online:a}}))}
                              />
                              <${De}
                                label="Sincronizar histórico"
                                description="Envia HistorySync recebido para webhooks cadastrados."
                                checked=${j.form.history_sync_enabled}
                                onChange=${a=>Q(u=>({...u,form:{...u.form,history_sync_enabled:a}}))}
                              />
                              <${De}
                                label="Usar proxy nesta sessão"
                                checked=${j.form.proxy_enabled}
                                onChange=${a=>Q(u=>({...u,form:{...u.form,proxy_enabled:a}}))}
                              />
                            </div>

                            <${oo}
                              label="Mensagem de rejeição"
                              value=${j.form.reject_calls_message}
                              rows="3"
                              placeholder="Opcional"
                              onInput=${a=>Q(u=>({...u,form:{...u.form,reject_calls_message:a}}))}
                            />
                            <${re}
                              label="Intervalo always-online (segundos)"
                              value=${j.form.always_online_interval_seconds}
                              type="number"
                              min="1"
                              step="1"
                              onInput=${a=>Q(u=>({...u,form:{...u.form,always_online_interval_seconds:a}}))}
                            />
                            <${re}
                              label="Dias de histórico"
                              value=${j.form.history_sync_days}
                              type="number"
                              min="1"
                              step="1"
                              onInput=${a=>Q(u=>({...u,form:{...u.form,history_sync_days:a}}))}
                            />

                            <div className="form-grid">
                              <${re}
                                label="Proxy protocolo"
                                value=${j.form.proxy_protocol}
                                placeholder="http, https ou socks5"
                                onInput=${a=>Q(u=>({...u,form:{...u.form,proxy_protocol:a}}))}
                              />
                              <${re}
                                label="Proxy host"
                                value=${j.form.proxy_host}
                                placeholder="proxy.exemplo.com"
                                onInput=${a=>Q(u=>({...u,form:{...u.form,proxy_host:a}}))}
                              />
                              <${re}
                                label="Proxy porta"
                                value=${j.form.proxy_port}
                                placeholder="8080"
                                onInput=${a=>Q(u=>({...u,form:{...u.form,proxy_port:a}}))}
                              />
                              <${re}
                                label="Proxy usuário"
                                value=${j.form.proxy_username}
                                placeholder="Opcional"
                                onInput=${a=>Q(u=>({...u,form:{...u.form,proxy_username:a}}))}
                              />
                              <${re}
                                label="Proxy senha"
                                value=${j.form.proxy_password}
                                type="password"
                                placeholder="Opcional"
                                onInput=${a=>Q(u=>({...u,form:{...u.form,proxy_password:a}}))}
                              />
                            </div>

                            <div className="inline-actions inline-actions--wrap">
                              <button
                                type="submit"
                                className="solid-button"
                                disabled=${j.loading||j.saving}
                              >
                                ${j.saving?"Salvando...":"Salvar configuração"}
                              </button>
                              <button
                                type="button"
                                className="ghost-button"
                                disabled=${j.loading}
                                onClick=${()=>ee(C)}
                              >
                                Recarregar
                              </button>
                              <button
                                type="button"
                                className="ghost-button ghost-button--danger"
                                disabled=${j.saving}
                                onClick=${Gr}
                              >
                                Limpar override
                              </button>
                            </div>
                          </form>
                        `:null}

                    ${ye==="webhooks"?K`
                          <form className="stack-form" onSubmit=${sr}>
                            <div className="form-grid">
                              <${re}
                                label="Webhook ID"
                                value=${$.form.id}
                                placeholder="vazio = auto"
                                onInput=${a=>f(u=>({...u,form:{...u.form,id:a}}))}
                              />
                              <${re}
                                label="Nome"
                                value=${$.form.name}
                                placeholder="Chatwoot principal"
                                onInput=${a=>f(u=>({...u,form:{...u.form,name:a}}))}
                              />
                            </div>

                            <${re}
                              label="URL base"
                              value=${$.form.url}
                              placeholder="https://seu-app.com/webhooks/whatsapp"
                              onInput=${a=>f(u=>({...u,form:{...u.form,url:a}}))}
                            />
                            <${re}
                              label="URL absoluta"
                              value=${$.form.url_absolute}
                              placeholder="https://seu-app.com/webhooks/whatsapp/5562..."
                              onInput=${a=>f(u=>({...u,form:{...u.form,url_absolute:a}}))}
                            />

                            <div className="form-grid">
                              <${re}
                                label="Header auth"
                                value=${$.form.header}
                                placeholder="Authorization"
                                onInput=${a=>f(u=>({...u,form:{...u.form,header:a}}))}
                              />
                              <${re}
                                label="Token"
                                value=${$.form.token}
                                placeholder="Bearer xxxxx"
                                onInput=${a=>f(u=>({...u,form:{...u.form,token:a}}))}
                              />
                            </div>

                            <${re}
                              label="Timeout (ms)"
                              value=${$.form.timeout_ms}
                              type="number"
                              min="1"
                              step="100"
                              onInput=${a=>f(u=>({...u,form:{...u.form,timeout_ms:a}}))}
                            />

                            <div className="check-grid">
                              <${De}
                                label="Webhook habilitado"
                                checked=${$.form.enabled}
                                onChange=${a=>f(u=>({...u,form:{...u.form,enabled:a}}))}
                              />
                              <${De}
                                label="Adicionar sessão à URL"
                                checked=${$.form.append_session}
                                onChange=${a=>f(u=>({...u,form:{...u.form,append_session:a}}))}
                              />
                              <${De}
                                label="Enviar incoming"
                                checked=${$.form.send_incoming_messages}
                                onChange=${a=>f(u=>({...u,form:{...u.form,send_incoming_messages:a}}))}
                              />
                              <${De}
                                label="Enviar outgoing"
                                checked=${$.form.send_outgoing_messages}
                                onChange=${a=>f(u=>({...u,form:{...u.form,send_outgoing_messages:a}}))}
                              />
                              <${De}
                                label="Enviar updates"
                                checked=${$.form.send_update_messages}
                                onChange=${a=>f(u=>({...u,form:{...u.form,send_update_messages:a}}))}
                              />
                              <${De}
                                label="Enviar grupos"
                                checked=${$.form.send_group_messages}
                                onChange=${a=>f(u=>({...u,form:{...u.form,send_group_messages:a}}))}
                              />
                              <${De}
                                label="Enviar newsletters"
                                checked=${$.form.send_newsletter_messages}
                                onChange=${a=>f(u=>({...u,form:{...u.form,send_newsletter_messages:a}}))}
                              />
                            </div>

                            <div className="inline-actions inline-actions--wrap">
                              <button
                                type="submit"
                                className="solid-button"
                                disabled=${$.loading||$.saving}
                              >
                                ${$.saving?"Salvando...":"Salvar webhook"}
                              </button>
                              <button
                                type="button"
                                className="ghost-button"
                                onClick=${()=>f(a=>({...a,form:dn()}))}
                              >
                                Limpar formulário
                              </button>
                              <button
                                type="button"
                                className="ghost-button"
                                disabled=${$.loading}
                                onClick=${()=>ce(C,{resetForm:!1})}
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
                              <span className="panel__hint">${$.items.length} item(ns)</span>
                            </div>

                            ${$.loading?K`
                                  <div className="empty-state empty-state--compact">
                                    <strong>Carregando webhooks…</strong>
                                  </div>
                                `:$.items.length===0?K`
                                    <div className="empty-state empty-state--compact">
                                      <strong>Nenhum webhook cadastrado</strong>
                                      <p>Preencha o formulário acima para salvar o primeiro callback.</p>
                                    </div>
                                  `:K`
                                    <div className="webhook-list">
                                      ${$.items.map(a=>K`
                                          <${_f}
                                            key=${a.id}
                                            item=${a}
                                            onEdit=${()=>f(u=>({...u,form:cf(a)}))}
                                            onDelete=${()=>In(a.id)}
                                          />
                                        `)}
                                    </div>
                                  `}
                          </div>
                        `:null}

                    ${ye==="templates"?K`
                          <form className="stack-form" onSubmit=${Xr}>
                            <div className="form-grid">
                              <${re}
                                label="Template ID"
                                value=${y.form.id}
                                placeholder="vazio = nome__idioma"
                                onInput=${a=>D(u=>({...u,form:{...u.form,id:a}}))}
                              />
                              <${re}
                                label="Nome Cloud API"
                                value=${y.form.name}
                                placeholder="boleto_disponivel"
                                onInput=${a=>D(u=>({...u,form:{...u.form,name:a}}))}
                              />
                              <${re}
                                label="Idioma"
                                value=${y.form.language}
                                placeholder="pt_BR"
                                onInput=${a=>D(u=>({...u,form:{...u.form,language:a}}))}
                              />
                              <${re}
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
                              <${re}
                                label="Texto do cabeçalho"
                                value=${y.form.header_text}
                                placeholder="Olá {{nome}}"
                                onInput=${a=>D(u=>({...u,form:{...u.form,header_text:a}}))}
                              />
                            </div>

                            ${y.form.header_type==="image"||y.form.header_type==="document"?K`
                                  <div className="form-grid">
                                    <${re}
                                      label="Link padrão da mídia"
                                      value=${y.form.header_link}
                                      placeholder="https://..."
                                      onInput=${a=>D(u=>({...u,form:{...u.form,header_link:a}}))}
                                    />
                                    <${re}
                                      label="Media ID padrão"
                                      value=${y.form.header_id}
                                      placeholder="opcional"
                                      onInput=${a=>D(u=>({...u,form:{...u.form,header_id:a}}))}
                                    />
                                    <${re}
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
                                onClick=${()=>je(C,{resetForm:!1})}
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

                            ${y.loading?K`
                                  <div className="empty-state empty-state--compact">
                                    <strong>Carregando templates…</strong>
                                  </div>
                                `:y.items.length===0?K`
                                    <div className="empty-state empty-state--compact">
                                      <strong>Nenhum template cadastrado</strong>
                                      <p>Preencha o formulário acima para habilitar envios type=template.</p>
                                    </div>
                                  `:K`
                                    <div className="webhook-list">
                                      ${y.items.map(a=>K`
                                          <${wf}
                                            key=${a.id}
                                            item=${a}
                                            onEdit=${()=>D(u=>({...u,form:df(a)}))}
                                            onDelete=${()=>gn(a.id)}
                                          />
                                        `)}
                                    </div>
                                  `}
                          </div>
                        `:null}
                  `:K`
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
        ${Fe.map(a=>K`
            <div key=${a.id} className=${`toast toast--${a.tone}`}>
              <strong>${a.tone==="error"?"Erro":a.tone==="success"?"OK":"Info"}</strong>
              <span>${a.message}</span>
            </div>
          `)}
      </div>

      <${yf} sessionId=${gt} onClose=${()=>yt("")} />
    </div>
  `}const Yu=document.getElementById("root");Yu&&ef.createRoot(Yu).render(K`<${kf} />`);
