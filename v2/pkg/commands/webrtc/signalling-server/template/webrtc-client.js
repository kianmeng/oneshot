"use strict";(()=>{var g=16384;var t=1*g;var e=8*g;var y="boundary";async function s(e,t,n){const o=(n==null?void 0:n.method)||"GET";let r=`${o} ${t} HTTP/1.1
`;let s=(n==null?void 0:n.headers)?n.headers:new Headers;if(s instanceof Headers){if(!s.has("User-Agent")){s.append("User-Agent",navigator.userAgent)}s.forEach((e,t)=>{r+=`${t}: ${e}
`})}else if(typeof s==="object"){if(s instanceof Array){var i=false;for(var a=0;a<s.length;a++){r+=`${s[a][0]}: ${s[a][1]}
`;if(s[a][0]==="User-Agent"){i=true}}if(!i){r+=`User-Agent: ${navigator.userAgent}
`}}else{if(!s["User-Agent"]){s["User-Agent"]=navigator.userAgent}for(const l in s){r+=`${l}: ${s[l]}
`}}}r+="\n";console.log("writing header: ",r);const c=d(e,r);try{c()}catch(e){return Promise.reject(e)}return Promise.resolve()}function d(n,e){var o=g;const r=function(){while(e.length){if(n.bufferedAmount>n.bufferedAmountLowThreshold){n.onbufferedamountlow=()=>{n.onbufferedamountlow=null;r()}}if(e.length<o){o=e.length}const t=e.slice(0,o);e=e.slice(o);try{n.send(t)}catch(e){if(e instanceof DOMException&&e.name==="InvalidStateError"){setTimeout(()=>{try{n.send(t)}catch(e){throw e}},500)}else{throw e}}if(o!=g){return}}};return r}async function i(t,e){let n=void 0;let o=void 0;let r=new Promise((e,t)=>{o=e;n=t});if(!e){try{t.send(new Uint8Array(0));t.send("");o()}catch(e){if(e instanceof DOMException&&e.name==="InvalidStateError"){setTimeout(()=>{try{t.send(new Uint8Array(0));t.send("")}catch(e){n(e)}},1e3)}else{n(e)}}return r}var s=void 0;if(e instanceof ReadableStream){const a=await e.getReader().read();e=a.value}else{if(e instanceof FormData){l(t,e).then(()=>{o()});return r}else if(e instanceof Blob){s=await e.arrayBuffer()}else if(e instanceof URLSearchParams){s=(new TextEncoder).encode(e.toString())}else if(typeof e==="string"){s=(new TextEncoder).encode(e)}else if(e instanceof ArrayBuffer){s=e}}const i=c(t,s,o);try{i()}catch(e){if(e instanceof DOMException&&e.name==="InvalidStateError"){setTimeout(()=>{try{i()}catch(e){n(e)}},500)}else{n(e)}}return r}function c(t,n,o){var r=g;const s=function(){while(n.byteLength){if(t.bufferedAmount>t.bufferedAmountLowThreshold){t.onbufferedamountlow=()=>{t.onbufferedamountlow=null;s()}}if(n.byteLength<r){r=n.byteLength}const e=n.slice(0,r);n=n.slice(r);t.send(e);if(r!=g){t.send("");if(o)o();return}}};return s}async function l(h,u){const m=new TextEncoder;return new Promise(async(e,t)=>{for(const i of u.entries()){var n=`--${y}
`;const a=i[0];const c=i[1];if(typeof c==="string"){n+=`Content-Disposition: form-data; name="${a}"

`;h.send(m.encode(n));h.send(m.encode(c))}else{const l=c;n+=`Content-Disposition: form-data; name="${a}"; filename="${l.name}"
`;if(l.type){n+=`Content-Type: ${l.type}

`}else{n+="Content-Type: application/octet-stream\n\n"}h.send(m.encode(n));var o;var r=new Promise((e,t)=>{o=e});const d=new FileReader;var s=0;d.onerror=e=>{console.log("Error reading file",e)};d.onabort=e=>{console.log("File reading aborted",e)};d.onload=e=>{const t=e.target.result;h.send(t);s+=t.byteLength;if(s<l.size){f(s)}else{o()}};const f=e=>{d.readAsArrayBuffer(l.slice(e,e+g))};f(s);await r}}h.send(m.encode(`
--${y}--
`));h.send("");e()})}function v(e){const t=e.split("\n");const n={};for(const o of t){const r=o.search(":");if(r===-1){continue}const s=o.slice(0,r);const i=o.slice(r+1);n[s]=i.trim()}return n}function b(e){if(!e.startsWith("HTTP/1.1")){throw new Error(`unexpected status line: ${e}`)}const t=e.split(" ");if(t.length<3){throw new Error(`unexpected status line: ${e}`)}const n=parseInt(t[1]);const o=t.slice(2).join(" ");return{status:n,statusText:o}}function o(w,r){w.bufferedAmountLowThreshold=t;w.binaryType="arraybuffer";let e=(e,t,l)=>{var d=()=>{};var n=()=>{};const o=new Promise((e,t)=>{d=e;n=t});var f="";var h=-1;var u="";var m={};const g=new MessageChannel;var p=[];w.onmessage=o=>{if(o.data instanceof ArrayBuffer){if(h===-1){const r=f.slice(0,f.search("\n"));f=f.slice(f.search("\n")+1);const s=b(r);h=s.status;u=s.statusText;m=v(f);f="";p.push(o.data);g.port1.postMessage(null);const i=new Headers;for(const a in m){i.append(a,m[a])}if(i.has("Content-Length")&&l){const c=parseInt(i.get("Content-Length"));l(-1,c)}let e={status:h,statusText:u,headers:i};let t=new ReadableStream({type:"bytes",start(e){},pull(r){return new Promise((o,e)=>{g.port2.onmessage=e=>{const t=p.shift();if(!t){w.send("");r.close();o();l==null?void 0:l(0);return}const n=t.byteLength;r.enqueue(new Uint8Array(t));o();l==null?void 0:l(n)}})}});let n=new Response(t,e);d(n)}else{const e=o.data;if(0<e.byteLength){p.push(e);g.port1.postMessage(null)}}}else if(typeof o.data==="string"){if(h===-1){f+=o.data}else{g.port1.postMessage(null)}}};if(!t){t={}}if(!t.headers){t.headers=new Headers}if((t==null?void 0:t.body)instanceof FormData){if(t.headers instanceof Headers){t.headers.append("Content-Type","multipart/form-data; boundary="+y)}else if(typeof t.headers==="object"){if(t.headers instanceof Array){t.headers.push(["Content-Type","multipart/form-data; boundary="+y])}else{t.headers["Content-Type"]="multipart/form-data; boundary="+y}}}if(r){if(t.headers instanceof Headers){t.headers.append("X-HTTPOverWebRTC-Authorization",r)}else if(typeof t.headers==="object"){if(t.headers instanceof Array){t.headers.push(["X-HTTPOverWebRTC-Authorization",r])}else{t.headers["X-HTTPOverWebRTC-Authorization"]=r}}}s(w,e,t).then(()=>{i(w,t==null?void 0:t.body).catch(e=>{console.log("rejecting promise 1: ",e);n(e)})}).catch(e=>{console.log("rejecting promise 2: ",e);n(e)});return o};return e}function w(e){var t;if(e instanceof HTMLScriptElement){(t=e.parentNode)==null?void 0:t.replaceChild(r(e),e)}else{var n=-1,o=e.childNodes;while(++n<o.length){w(o[n])}}}function r(e){var t=document.createElement("script");t.text=e.innerHTML;var n=-1,o=e.attributes,r;while(++n<o.length){t.setAttribute((r=o[n]).name,r.value)}return t}function T(e,t){const n=document.createElement("a");n.setAttribute("style","display: none");document.body.appendChild(n);const o=new Blob([e],{type:"stream/octet"});const r=window.URL.createObjectURL(o);n.href=r;n.download=t;n.click();window.URL.revokeObjectURL(r)}var P=/^text\/.*$/;async function n(e,t,n=fetch){const o=document.createElement("span");document.body.innerHTML="";document.body.appendChild(o);var r=0;var s=0;const i=(e,t)=>{if(e===-1&&t){r=t}else if(0<e){s+=e;if(r){const n=(s/r*100).toFixed(2);o.innerText=`receiving data: ${n}%`}else{o.innerText=`receiving data: ${s} bytes`}}else{o.innerText=`receiving data: done`}return Promise.resolve()};const a=n;var c=await a(e,t,i);const l=c.headers;var d=l.get("Content-Type")?l.get("Content-Type"):"";d=d.split(";")[0];var f=l.get("Content-Disposition")?l.get("Content-Disposition"):"";if(!d){d="text/plain"}const h=A(f);if(h){const u=await c.blob();T(u,h);return}if(d.match(P)){const m=await c.text();if(d==="text/html"){const g=new DOMParser;const p=g.parseFromString(m,"text/html");document.body=p.body;w(document.body)}else{document.body.innerText=m;document.body.innerHTML=`<pre>${m}</pre>`}}else if(d.startsWith("application/")){const m=await c.blob();let e=new Blob([m],{type:d});let t=URL.createObjectURL(e);window.open(t,"_self")}else{console.log(`falling back to displaying body as preformatted text`);const m=await c.text();document.body.innerText=m;document.body.innerHTML=`<pre>${m}</pre>`}}function A(e){if(!e||!e.includes("attachment")){return""}const t=/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/;const n=t.exec(e);if(n!=null&&n[1]){return n[1].replace(/['"]/g,"")}return""}var a=class{constructor(e){this.answered=false;this.connectionPromiseResolve=()=>{};this.connectionPromiseReject=()=>{};this.resolveAnswerPromise=()=>{};this.rejectAnswerPromise=()=>{};this.answerPromise=new Promise((e,t)=>{this.resolveAnswerPromise=e;this.rejectAnswerPromise=t});this.peerConnection=new RTCPeerConnection(e);this.connectionPromise=new Promise((e,t)=>{this.connectionPromiseResolve=e;this.connectionPromiseReject=t});this._configurePeerConnection()}_configurePeerConnection(){const n=this.peerConnection;n.onicegatheringstatechange=e=>{console.log("onicegatheringstatechange",n.iceGatheringState);console.log("event",e);let t=e.target;console.log("target",t);if(t.iceGatheringState==="complete"&&t.localDescription){this.resolveAnswerPromise(t.localDescription)}};n.ondatachannel=e=>{console.log("ondatachannel",e);this._fetch=o(e.channel,this.baToken);window.fetch=this._fetch;window.rtcReady=true;this.connectionPromiseResolve()};n.onnegotiationneeded=e=>{console.log("onnegotiationneeded")};n.onsignalingstatechange=e=>{console.log("onsignalingstatechange",n.signalingState)};n.oniceconnectionstatechange=e=>{console.log("oniceconnectionstatechange",e);console.log("oniceconnectionstatechange",n.iceConnectionState)};n.onicecandidate=e=>{console.log("onicecandidate",e)}}answerOffer(e){if(this.answered){return{Answer:Promise.reject("already answered"),ConnectionEstablished:Promise.reject("already answered")}}let t=e.sdp.match(/a=BasicAuthToken:(.*)/);if(t&&t.length>1){this.baToken=t[1]}const n=this.peerConnection;const o=e=>{n.setLocalDescription(e).then(()=>{this.answered=true;this.resolveAnswerPromise(new RTCSessionDescription(e))})};n.setRemoteDescription(e).then(()=>n.createAnswer().then(o)).catch(e=>{console.error(e);this.rejectAnswerPromise(e);this.connectionPromiseReject(e)});return{Answer:this.answerPromise,ConnectionEstablished:this.connectionPromise}}fetch(e,t){if(!this._fetch){return Promise.reject("HTTPOverWebRTC fetch not ready")}return this._fetch(e,t)}visit(e,t){if(!this._fetch){return Promise.reject("HTTPOverWebRTC fetch not ready")}return n(e,t,this._fetch.bind(this))}connected(){return this._fetch!==void 0}};window.WebRTCClient=a})();