import{d as ze,A as u,W as g,X as R,B as v,b as L,C as j,e as m,f as y,F as G,n as I,g as U,E as _,u as ce,I as se,c as b,a as ne,y as oe,z as ae,r as O,s as uo,b5 as vo,t as mo,h as W,k as te,ar as ho,ah as Me,T as he,ai as fo,b0 as Pe,l as $,bm as Oe,m as Z,aP as Ke,bn as Ve,al as De,G as X,$ as V,Y as le,ao as po,ak as Ie,aj as re,aU as Se,q as be,an as go,aX as fe,V as bo,as as Ne,aS as Co,J as xo,O as Y,Q as M,N as D,a7 as Be,S as $e,a6 as yo,bo as zo,bp as Io,ba as So}from"./index-C-2PX-C-.js";import{n as wo,t as Ao,C as Ro,T as Ho,D as Ue,V as ko,m as _o,c as pe}from"./Dropdown-CM88cknF.js";import{f as ge,u as Ce}from"./format-length-icqZkaSq.js";import{_ as To}from"./_plugin-vue_export-helper-DlAUqK2U.js";import"./next-frame-once-C5Ksf8W7.js";var Po={fontWeightActive:"400"};function No(e){const{fontSize:t,textColor3:o,textColor2:r,borderRadius:a,buttonColor2Hover:i,buttonColor2Pressed:s}=e;return{...Po,fontSize:t,itemLineHeight:"1.25",itemTextColor:o,itemTextColorHover:r,itemTextColorPressed:r,itemTextColorActive:r,itemBorderRadius:a,itemColorHover:i,itemColorPressed:s,separatorColor:o}}const Bo={common:ze,self:No};var $o=u("breadcrumb",`
 white-space: nowrap;
 cursor: default;
 line-height: var(--n-item-line-height);
`,[g("ul",`
 list-style: none;
 padding: 0;
 margin: 0;
 `),g("a",`
 color: inherit;
 text-decoration: inherit;
 `),u("breadcrumb-item",`
 font-size: var(--n-font-size);
 transition: color .3s var(--n-bezier);
 display: inline-flex;
 align-items: center;
 `,[u("icon",`
 font-size: 18px;
 vertical-align: -.2em;
 transition: color .3s var(--n-bezier);
 color: var(--n-item-text-color);
 `),g("&:not(:last-child)",[R("clickable",[v("link",`
 cursor: pointer;
 `,[g("&:hover",`
 background-color: var(--n-item-color-hover);
 `),g("&:active",`
 background-color: var(--n-item-color-pressed); 
 `)])])]),v("link",`
 padding: 4px;
 border-radius: var(--n-item-border-radius);
 transition:
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 color: var(--n-item-text-color);
 position: relative;
 `,[g("&:hover",`
 color: var(--n-item-text-color-hover);
 `,[u("icon",`
 color: var(--n-item-text-color-hover);
 `)]),g("&:active",`
 color: var(--n-item-text-color-pressed);
 `,[u("icon",`
 color: var(--n-item-text-color-pressed);
 `)])]),v("separator",`
 margin: 0 8px;
 color: var(--n-separator-color);
 transition: color .3s var(--n-bezier);
 user-select: none;
 -webkit-user-select: none;
 `),g("&:last-child",[v("link",`
 font-weight: var(--n-font-weight-active);
 cursor: unset;
 color: var(--n-item-text-color-active);
 `,[u("icon",`
 color: var(--n-item-text-color-active);
 `)]),v("separator",`
 display: none;
 `)])])]);const je=ne("n-breadcrumb"),Eo={...j.props,separator:{type:String,default:"/"}};var Lo=L({name:"Breadcrumb",props:Eo,setup(e){const{mergedClsPrefixRef:t,inlineThemeDisabled:o}=ce(e),r=j("Breadcrumb","-breadcrumb",$o,Bo,e,t);oe(je,{separatorRef:ae(e,"separator"),mergedClsPrefixRef:t});const a=b(()=>{const{common:{cubicBezierEaseInOut:s},self:{separatorColor:d,itemTextColor:c,itemTextColorHover:h,itemTextColorPressed:H,itemTextColorActive:z,fontSize:p,fontWeightActive:S,itemBorderRadius:P,itemColorHover:x,itemColorPressed:T,itemLineHeight:B}}=r.value;return{"--n-font-size":p,"--n-bezier":s,"--n-item-text-color":c,"--n-item-text-color-hover":h,"--n-item-text-color-pressed":H,"--n-item-text-color-active":z,"--n-separator-color":d,"--n-item-color-hover":x,"--n-item-color-pressed":T,"--n-item-border-radius":P,"--n-font-weight-active":S,"--n-item-line-height":B}}),i=o?se("breadcrumb",void 0,a,e):void 0;return{mergedClsPrefix:t,cssVars:o?void 0:a,themeClass:i==null?void 0:i.themeClass,onRender:i==null?void 0:i.onRender}},render(){var e;return(e=this.onRender)==null||e.call(this),m(),y("nav",{class:_([`${this.mergedClsPrefix}-breadcrumb`,this.themeClass]),style:U(this.cssVars),"aria-label":"Breadcrumb"},[G("ul",null,[I(()=>{var t,o;return(o=(t=this.$slots).default)==null?void 0:o.call(t)})])],6)}});function Fo(e=mo?window:null){const t=()=>{const{hash:a,host:i,hostname:s,href:d,origin:c,pathname:h,port:H,protocol:z,search:p}=(e==null?void 0:e.location)||{};return{hash:a,host:i,hostname:s,href:d,origin:c,pathname:h,port:H,protocol:z,search:p}},o=O(t()),r=()=>{o.value=t()};return uo(()=>{e&&(e.addEventListener("popstate",r),e.addEventListener("hashchange",r))}),vo(()=>{e&&(e.removeEventListener("popstate",r),e.removeEventListener("hashchange",r))}),o}const Mo={separator:String,href:String,clickable:{type:Boolean,default:!0},showSeparator:{type:Boolean,default:!0},onClick:Function};var Oo=L({name:"BreadcrumbItem",props:Mo,slots:Object,setup(e,{slots:t}){const o=W(je,null);if(!o)return()=>null;const{separatorRef:r,mergedClsPrefixRef:a}=o,i=Fo(),s=b(()=>e.href?"a":"span"),d=b(()=>i.value.href===e.href?"location":null);return()=>{const{value:c}=a;return m(),y("li",{class:_([`${c}-breadcrumb-item`,e.clickable&&`${c}-breadcrumb-item--clickable`])},[I(()=>te(s.value,{class:`${c}-breadcrumb-item__link`,"aria-current":d.value,href:e.href,onClick:e.onClick},t)),I(()=>e.showSeparator&&(m(),y("span",{class:_(`${c}-breadcrumb-item__separator`),"aria-hidden":"true"},[I(()=>ho(t.separator,()=>[e.separator??r.value]))],2)))],2)}}});function Ko(e,t,o,r){return{itemColorHoverInverted:"#0000",itemColorActiveInverted:t,itemColorActiveHoverInverted:t,itemColorActiveCollapsedInverted:t,itemTextColorInverted:e,itemTextColorHoverInverted:o,itemTextColorChildActiveInverted:o,itemTextColorChildActiveHoverInverted:o,itemTextColorActiveInverted:o,itemTextColorActiveHoverInverted:o,itemTextColorHorizontalInverted:e,itemTextColorHoverHorizontalInverted:o,itemTextColorChildActiveHorizontalInverted:o,itemTextColorChildActiveHoverHorizontalInverted:o,itemTextColorActiveHorizontalInverted:o,itemTextColorActiveHoverHorizontalInverted:o,itemIconColorInverted:e,itemIconColorHoverInverted:o,itemIconColorActiveInverted:o,itemIconColorActiveHoverInverted:o,itemIconColorChildActiveInverted:o,itemIconColorChildActiveHoverInverted:o,itemIconColorCollapsedInverted:e,itemIconColorHorizontalInverted:e,itemIconColorHoverHorizontalInverted:o,itemIconColorActiveHorizontalInverted:o,itemIconColorActiveHoverHorizontalInverted:o,itemIconColorChildActiveHorizontalInverted:o,itemIconColorChildActiveHoverHorizontalInverted:o,arrowColorInverted:e,arrowColorHoverInverted:o,arrowColorActiveInverted:o,arrowColorActiveHoverInverted:o,arrowColorChildActiveInverted:o,arrowColorChildActiveHoverInverted:o,groupTextColorInverted:r}}function Vo(e){const{borderRadius:t,textColor3:o,primaryColor:r,textColor2:a,textColor1:i,fontSize:s,dividerColor:d,hoverColor:c,primaryColorHover:h}=e;return{borderRadius:t,color:"#0000",groupTextColor:o,itemColorHover:c,itemColorActive:he(r,{alpha:.1}),itemColorActiveHover:he(r,{alpha:.1}),itemColorActiveCollapsed:he(r,{alpha:.1}),itemTextColor:a,itemTextColorHover:a,itemTextColorActive:r,itemTextColorActiveHover:r,itemTextColorChildActive:r,itemTextColorChildActiveHover:r,itemTextColorHorizontal:a,itemTextColorHoverHorizontal:h,itemTextColorActiveHorizontal:r,itemTextColorActiveHoverHorizontal:r,itemTextColorChildActiveHorizontal:r,itemTextColorChildActiveHoverHorizontal:r,itemIconColor:i,itemIconColorHover:i,itemIconColorActive:r,itemIconColorActiveHover:r,itemIconColorChildActive:r,itemIconColorChildActiveHover:r,itemIconColorCollapsed:i,itemIconColorHorizontal:i,itemIconColorHoverHorizontal:h,itemIconColorActiveHorizontal:r,itemIconColorActiveHoverHorizontal:r,itemIconColorChildActiveHorizontal:r,itemIconColorChildActiveHoverHorizontal:r,itemHeight:"42px",arrowColor:a,arrowColorHover:a,arrowColorActive:r,arrowColorActiveHover:r,arrowColorChildActive:r,arrowColorChildActiveHover:r,colorInverted:"#0000",borderColorHorizontal:"#0000",fontSize:s,dividerColor:d,...Ko("#BBB",r,"#FFF","#AAA")}}const Do=Me({name:"Menu",common:ze,peers:{Tooltip:Ao,Dropdown:wo},self:Vo});function Uo(e){const{baseColor:t,textColor2:o,bodyColor:r,cardColor:a,dividerColor:i,actionColor:s,scrollbarColor:d,scrollbarColorHover:c,invertedColor:h}=e;return{textColor:o,textColorInverted:"#FFF",color:r,colorEmbedded:s,headerColor:a,headerColorInverted:h,footerColor:s,footerColorInverted:h,headerBorderColor:i,headerBorderColorInverted:h,footerBorderColor:i,footerBorderColorInverted:h,siderBorderColor:i,siderBorderColorInverted:h,siderColor:a,siderColorInverted:h,siderToggleButtonBorder:`1px solid ${i}`,siderToggleButtonColor:t,siderToggleButtonIconColor:o,siderToggleButtonIconColorInverted:o,siderToggleBarColor:Pe(r,d),siderToggleBarColorHover:Pe(r,c),__invertScrollbar:"true"}}const we=Me({name:"Layout",common:ze,peers:{Scrollbar:fo},self:Uo}),Ge=ne("n-layout-sider"),Ae={type:String,default:"static"};var jo=u("layout",`
 color: var(--n-text-color);
 background-color: var(--n-color);
 box-sizing: border-box;
 position: relative;
 z-index: auto;
 flex: auto;
 overflow: hidden;
 transition:
 box-shadow .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
`,[u("layout-scroll-container",`
 overflow-x: hidden;
 box-sizing: border-box;
 height: 100%;
 `),R("absolute-positioned",`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]);const Go={embedded:Boolean,position:Ae,nativeScrollbar:{type:Boolean,default:!0},scrollbarProps:Object,onScroll:Function,contentClass:String,contentStyle:{type:[String,Object],default:""},hasSider:Boolean,siderPlacement:{type:String,default:"left"}},qe=ne("n-layout");function Ye(e){return L({name:e?"LayoutContent":"Layout",props:{...j.props,...Go},setup(t){const o=O(null),r=O(null),{mergedClsPrefixRef:a,inlineThemeDisabled:i}=ce(t),s=j("Layout","-layout",jo,we,t,a);function d(x,T){if(t.nativeScrollbar){const{value:B}=o;B&&(T===void 0?B.scrollTo(x):B.scrollTo(x,T))}else{const{value:B}=r;B&&B.scrollTo(x,T)}}oe(qe,t);let c=0,h=0;const H=x=>{var B;const T=x.target;c=T.scrollLeft,h=T.scrollTop,(B=t.onScroll)==null||B.call(t,x)};Ve(()=>{if(t.nativeScrollbar){const x=o.value;x&&(x.scrollTop=h,x.scrollLeft=c)}});const z={display:"flex",flexWrap:"nowrap",width:"100%",flexDirection:"row"},p={scrollTo:d},S=b(()=>{const{common:{cubicBezierEaseInOut:x},self:T}=s.value;return{"--n-bezier":x,"--n-color":t.embedded?T.colorEmbedded:T.color,"--n-text-color":T.textColor}}),P=i?se("layout",b(()=>t.embedded?"e":""),S,t):void 0;return{mergedClsPrefix:a,scrollableElRef:o,scrollbarInstRef:r,hasSiderStyle:z,mergedTheme:s,handleNativeElScroll:H,cssVars:i?void 0:S,themeClass:P==null?void 0:P.themeClass,onRender:P==null?void 0:P.onRender,...p}},render(){var i;const{mergedClsPrefix:t,hasSider:o}=this;(i=this.onRender)==null||i.call(this);const r=o?this.hasSiderStyle:void 0,a=[this.themeClass,e&&`${t}-layout-content`,`${t}-layout`,`${t}-layout--${this.position}-positioned`];return m(),y("div",{class:_(a),style:U(this.cssVars)},[this.nativeScrollbar?(m(),y("div",{key:0,ref:"scrollableElRef",class:_([`${t}-layout-scroll-container`,this.contentClass]),style:U([this.contentStyle,r]),onScroll:this.handleNativeElScroll},[I(()=>{var s,d;return(d=(s=this.$slots).default)==null?void 0:d.call(s)})],46,["onScroll"])):(m(),$(Ke,Z({key:1},this.scrollbarProps,{onScroll:this.onScroll,ref:"scrollbarInstRef",theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar,contentClass:this.contentClass,contentStyle:[this.contentStyle,r]}),Oe(this.$slots),1040,["onScroll","theme","themeOverrides","contentClass","contentStyle"]))],6)}})}var Ee=Ye(!1),qo=Ye(!0),Yo=u("layout-header",`
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 box-sizing: border-box;
 width: 100%;
 background-color: var(--n-color);
 color: var(--n-text-color);
`,[R("absolute-positioned",`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 `),R("bordered",`
 border-bottom: solid 1px var(--n-border-color);
 `)]);const Xo={position:Ae,inverted:Boolean,bordered:Boolean};var Wo=L({name:"LayoutHeader",props:{...j.props,...Xo},setup(e){const{mergedClsPrefixRef:t,inlineThemeDisabled:o}=ce(e),r=j("Layout","-layout-header",Yo,we,e,t),a=b(()=>{const{common:{cubicBezierEaseInOut:s},self:d}=r.value,c={"--n-bezier":s};return e.inverted?(c["--n-color"]=d.headerColorInverted,c["--n-text-color"]=d.textColorInverted,c["--n-border-color"]=d.headerBorderColorInverted):(c["--n-color"]=d.headerColor,c["--n-text-color"]=d.textColor,c["--n-border-color"]=d.headerBorderColor),c}),i=o?se("layout-header",b(()=>e.inverted?"a":"b"),a,e):void 0;return{mergedClsPrefix:t,cssVars:o?void 0:a,themeClass:i==null?void 0:i.themeClass,onRender:i==null?void 0:i.onRender}},render(){var t;const{mergedClsPrefix:e}=this;return(t=this.onRender)==null||t.call(this),m(),y("div",{class:_([`${e}-layout-header`,this.themeClass,this.position&&`${e}-layout-header--${this.position}-positioned`,this.bordered&&`${e}-layout-header--bordered`]),style:U(this.cssVars)},[I(()=>{var o,r;return(r=(o=this.$slots).default)==null?void 0:r.call(o)})],6)}}),Jo=u("layout-sider",`
 flex-shrink: 0;
 box-sizing: border-box;
 position: relative;
 z-index: 1;
 color: var(--n-text-color);
 transition:
 color .3s var(--n-bezier),
 border-color .3s var(--n-bezier),
 min-width .3s var(--n-bezier),
 max-width .3s var(--n-bezier),
 transform .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 background-color: var(--n-color);
 display: flex;
 justify-content: flex-end;
`,[R("bordered",[v("border",`
 content: "";
 position: absolute;
 top: 0;
 bottom: 0;
 width: 1px;
 background-color: var(--n-border-color);
 transition: background-color .3s var(--n-bezier);
 `)]),v("left-placement",[R("bordered",[v("border",`
 right: 0;
 `)])]),R("right-placement",`
 justify-content: flex-start;
 `,[R("bordered",[v("border",`
 left: 0;
 `)]),R("collapsed",[u("layout-toggle-button",[u("base-icon",`
 transform: rotate(180deg);
 `)]),u("layout-toggle-bar",[g("&:hover",[v("top",{transform:"rotate(-12deg) scale(1.15) translateY(-2px)"}),v("bottom",{transform:"rotate(12deg) scale(1.15) translateY(2px)"})])])]),u("layout-toggle-button",`
 left: 0;
 transform: translateX(-50%) translateY(-50%);
 `,[u("base-icon",`
 transform: rotate(0);
 `)]),u("layout-toggle-bar",`
 left: -28px;
 transform: rotate(180deg);
 `,[g("&:hover",[v("top",{transform:"rotate(12deg) scale(1.15) translateY(-2px)"}),v("bottom",{transform:"rotate(-12deg) scale(1.15) translateY(2px)"})])])]),R("collapsed",[u("layout-toggle-bar",[g("&:hover",[v("top",{transform:"rotate(-12deg) scale(1.15) translateY(-2px)"}),v("bottom",{transform:"rotate(12deg) scale(1.15) translateY(2px)"})])]),u("layout-toggle-button",[u("base-icon",`
 transform: rotate(0);
 `)])]),u("layout-toggle-button",`
 transition:
 color .3s var(--n-bezier),
 right .3s var(--n-bezier),
 left .3s var(--n-bezier),
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 cursor: pointer;
 width: 24px;
 height: 24px;
 position: absolute;
 top: 50%;
 right: 0;
 border-radius: 50%;
 display: flex;
 align-items: center;
 justify-content: center;
 font-size: 18px;
 color: var(--n-toggle-button-icon-color);
 border: var(--n-toggle-button-border);
 background-color: var(--n-toggle-button-color);
 box-shadow: 0 2px 4px 0px rgba(0, 0, 0, .06);
 transform: translateX(50%) translateY(-50%);
 z-index: 1;
 `,[u("base-icon",`
 transition: transform .3s var(--n-bezier);
 transform: rotate(180deg);
 `)]),u("layout-toggle-bar",`
 cursor: pointer;
 height: 72px;
 width: 32px;
 position: absolute;
 top: calc(50% - 36px);
 right: -28px;
 `,[v("top, bottom",`
 position: absolute;
 width: 4px;
 border-radius: 2px;
 height: 38px;
 left: 14px;
 transition: 
 background-color .3s var(--n-bezier),
 transform .3s var(--n-bezier);
 `),v("bottom",`
 position: absolute;
 top: 34px;
 `),g("&:hover",[v("top",{transform:"rotate(12deg) scale(1.15) translateY(-2px)"}),v("bottom",{transform:"rotate(-12deg) scale(1.15) translateY(2px)"})]),v("top, bottom",{backgroundColor:"var(--n-toggle-bar-color)"}),g("&:hover",[v("top, bottom",{backgroundColor:"var(--n-toggle-bar-color-hover)"})])]),v("border",`
 position: absolute;
 top: 0;
 right: 0;
 bottom: 0;
 width: 1px;
 transition: background-color .3s var(--n-bezier);
 `),u("layout-sider-scroll-container",`
 flex-grow: 1;
 flex-shrink: 0;
 box-sizing: border-box;
 height: 100%;
 opacity: 0;
 transition: opacity .3s var(--n-bezier);
 max-width: 100%;
 `),R("show-content",[u("layout-sider-scroll-container",{opacity:1})]),R("absolute-positioned",`
 position: absolute;
 left: 0;
 top: 0;
 bottom: 0;
 `)]);const Qo=["onClick"];var Zo=L({props:{clsPrefix:{type:String,required:!0},onClick:Function},render(){const{clsPrefix:e}=this;return m(),y("div",{onClick:this.onClick,class:_(`${e}-layout-toggle-bar`)},[G("div",{class:_(`${e}-layout-toggle-bar__top`)},null,2),G("div",{class:_(`${e}-layout-toggle-bar__bottom`)},null,2)],10,Qo)}});const et=["onClick"];var ot=L({name:"LayoutToggleButton",props:{clsPrefix:{type:String,required:!0},onClick:Function},render(){const{clsPrefix:e}=this;return m(),y("div",{class:_(`${e}-layout-toggle-button`),onClick:this.onClick},[(m(),$(De,{clsPrefix:e},{default:()=>(m(),$(Ro))},1032,["clsPrefix"]))],10,et)}});const tt=["onTransitionend"],rt={position:Ae,bordered:Boolean,collapsedWidth:{type:Number,default:48},width:{type:[Number,String],default:272},contentClass:String,contentStyle:{type:[String,Object],default:""},collapseMode:{type:String,default:"transform"},collapsed:{type:Boolean,default:void 0},defaultCollapsed:Boolean,showCollapsedContent:{type:Boolean,default:!0},showTrigger:{type:[Boolean,String],default:!1},nativeScrollbar:{type:Boolean,default:!0},inverted:Boolean,scrollbarProps:Object,triggerClass:String,triggerStyle:[String,Object],collapsedTriggerClass:String,collapsedTriggerStyle:[String,Object],"onUpdate:collapsed":[Function,Array],onUpdateCollapsed:[Function,Array],onAfterEnter:Function,onAfterLeave:Function,onExpand:[Function,Array],onCollapse:[Function,Array],onScroll:Function};var nt=L({name:"LayoutSider",props:{...j.props,...rt},setup(e){const t=W(qe),o=O(null),r=O(null),a=O(e.defaultCollapsed),i=Ce(ae(e,"collapsed"),a),s=b(()=>ge(i.value?e.collapsedWidth:e.width)),d=b(()=>e.collapseMode!=="transform"?{}:{minWidth:ge(e.width)}),c=b(()=>t?t.siderPlacement:"left");function h(N,A){if(e.nativeScrollbar){const{value:w}=o;w&&(A===void 0?w.scrollTo(N):w.scrollTo(N,A))}else{const{value:w}=r;w&&w.scrollTo(N,A)}}function H(){const{"onUpdate:collapsed":N,onUpdateCollapsed:A,onExpand:w,onCollapse:Q}=e,{value:q}=i;A&&V(A,!q),N&&V(N,!q),a.value=!q,q?w&&V(w):Q&&V(Q)}let z=0,p=0;const S=N=>{var w;const A=N.target;z=A.scrollLeft,p=A.scrollTop,(w=e.onScroll)==null||w.call(e,N)};Ve(()=>{if(e.nativeScrollbar){const N=o.value;N&&(N.scrollTop=p,N.scrollLeft=z)}}),oe(Ge,{collapsedRef:i,collapseModeRef:ae(e,"collapseMode")});const{mergedClsPrefixRef:P,inlineThemeDisabled:x}=ce(e),T=j("Layout","-layout-sider",Jo,we,e,P);function B(N){var A,w;N.propertyName==="max-width"&&(i.value?(A=e.onAfterLeave)==null||A.call(e):(w=e.onAfterEnter)==null||w.call(e))}const ie={scrollTo:h},J=b(()=>{const{common:{cubicBezierEaseInOut:N},self:A}=T.value,{siderToggleButtonColor:w,siderToggleButtonBorder:Q,siderToggleBarColor:q,siderToggleBarColorHover:me}=A,F={"--n-bezier":N,"--n-toggle-button-color":w,"--n-toggle-button-border":Q,"--n-toggle-bar-color":q,"--n-toggle-bar-color-hover":me};return e.inverted?(F["--n-color"]=A.siderColorInverted,F["--n-text-color"]=A.textColorInverted,F["--n-border-color"]=A.siderBorderColorInverted,F["--n-toggle-button-icon-color"]=A.siderToggleButtonIconColorInverted,F.__invertScrollbar=A.__invertScrollbar):(F["--n-color"]=A.siderColor,F["--n-text-color"]=A.textColor,F["--n-border-color"]=A.siderBorderColor,F["--n-toggle-button-icon-color"]=A.siderToggleButtonIconColor),F}),K=x?se("layout-sider",b(()=>e.inverted?"a":"b"),J,e):void 0;return{scrollableElRef:o,scrollbarInstRef:r,mergedClsPrefix:P,mergedTheme:T,styleMaxWidth:s,mergedCollapsed:i,scrollContainerStyle:d,siderPlacement:c,handleNativeElScroll:S,handleTransitionend:B,handleTriggerClick:H,inlineThemeDisabled:x,cssVars:J,themeClass:K==null?void 0:K.themeClass,onRender:K==null?void 0:K.onRender,...ie}},render(){var r;const{mergedClsPrefix:e,mergedCollapsed:t,showTrigger:o}=this;return(r=this.onRender)==null||r.call(this),m(),y("aside",{class:_([`${e}-layout-sider`,this.themeClass,`${e}-layout-sider--${this.position}-positioned`,`${e}-layout-sider--${this.siderPlacement}-placement`,this.bordered&&`${e}-layout-sider--bordered`,t&&`${e}-layout-sider--collapsed`,(!t||this.showCollapsedContent)&&`${e}-layout-sider--show-content`]),onTransitionend:this.handleTransitionend,style:U([this.inlineThemeDisabled?void 0:this.cssVars,{maxWidth:this.styleMaxWidth,width:ge(this.width)}])},[this.nativeScrollbar?(m(),y("div",{key:1,class:_([`${e}-layout-sider-scroll-container`,this.contentClass]),onScroll:this.handleNativeElScroll,style:U([this.scrollContainerStyle,{overflow:"auto"},this.contentStyle]),ref:"scrollableElRef"},[I(()=>{var a,i;return(i=(a=this.$slots).default)==null?void 0:i.call(a)})],46,["onScroll"])):(m(),$(Ke,Z({key:0},this.scrollbarProps,{onScroll:this.onScroll,ref:"scrollbarInstRef",style:this.scrollContainerStyle,contentStyle:this.contentStyle,contentClass:this.contentClass,theme:this.mergedTheme.peers.Scrollbar,themeOverrides:this.mergedTheme.peerOverrides.Scrollbar,builtinThemeOverrides:this.inverted&&this.cssVars.__invertScrollbar==="true"?{colorHover:"rgba(255, 255, 255, .4)",color:"rgba(255, 255, 255, .3)"}:void 0}),Oe(this.$slots),1040,["onScroll","style","contentStyle","contentClass","theme","themeOverrides","builtinThemeOverrides"])),o?(m(),y(X,{key:2},[o==="bar"?(m(),$(Zo,{key:0,clsPrefix:e,class:_(t?this.collapsedTriggerClass:this.triggerClass),style:U(t?this.collapsedTriggerStyle:this.triggerStyle),onClick:this.handleTriggerClick},null,8,["clsPrefix","class","style","onClick"])):(m(),$(ot,{key:1,clsPrefix:e,class:_(t?this.collapsedTriggerClass:this.triggerClass),style:U(t?this.collapsedTriggerStyle:this.triggerStyle),onClick:this.handleTriggerClick},null,8,["clsPrefix","class","style","onClick"]))],64)):I(()=>null),this.bordered?(m(),y("div",{key:4,class:_(`${e}-layout-sider__border`)},null,2)):I(()=>null)],46,tt)}});const de=ne("n-menu"),Xe=ne("n-submenu"),Re=ne("n-menu-item-group"),Le=[g("&::before","background-color: var(--n-item-color-hover);"),v("arrow",`
 color: var(--n-arrow-color-hover);
 `),v("icon",`
 color: var(--n-item-icon-color-hover);
 `),u("menu-item-content-header",`
 color: var(--n-item-text-color-hover);
 `,[g("a",`
 color: var(--n-item-text-color-hover);
 `),v("extra",`
 color: var(--n-item-text-color-hover);
 `)])],Fe=[v("icon",`
 color: var(--n-item-icon-color-hover-horizontal);
 `),u("menu-item-content-header",`
 color: var(--n-item-text-color-hover-horizontal);
 `,[g("a",`
 color: var(--n-item-text-color-hover-horizontal);
 `),v("extra",`
 color: var(--n-item-text-color-hover-horizontal);
 `)])];var it=g([u("menu",`
 background-color: var(--n-color);
 color: var(--n-item-text-color);
 overflow: hidden;
 transition: background-color .3s var(--n-bezier);
 box-sizing: border-box;
 font-size: var(--n-font-size);
 padding-bottom: 6px;
 `,[R("horizontal",`
 max-width: 100%;
 width: 100%;
 display: flex;
 overflow: hidden;
 padding-bottom: 0;
 `,[u("submenu","margin: 0;"),u("menu-item","margin: 0;"),u("menu-item-content",`
 padding: 0 20px;
 border-bottom: 2px solid #0000;
 `,[g("&::before","display: none;"),R("selected","border-bottom: 2px solid var(--n-border-color-horizontal)")]),u("menu-item-content",[R("selected",[v("icon","color: var(--n-item-icon-color-active-horizontal);"),u("menu-item-content-header",`
 color: var(--n-item-text-color-active-horizontal);
 `,[g("a","color: var(--n-item-text-color-active-horizontal);"),v("extra","color: var(--n-item-text-color-active-horizontal);")])]),R("child-active",`
 border-bottom: 2px solid var(--n-border-color-horizontal);
 `,[u("menu-item-content-header",`
 color: var(--n-item-text-color-child-active-horizontal);
 `,[g("a",`
 color: var(--n-item-text-color-child-active-horizontal);
 `),v("extra",`
 color: var(--n-item-text-color-child-active-horizontal);
 `)]),v("icon",`
 color: var(--n-item-icon-color-child-active-horizontal);
 `)]),le("disabled",[le("selected, child-active",[g("&:focus-within",Fe)]),R("selected",[ee(null,[v("icon","color: var(--n-item-icon-color-active-hover-horizontal);"),u("menu-item-content-header",`
 color: var(--n-item-text-color-active-hover-horizontal);
 `,[g("a","color: var(--n-item-text-color-active-hover-horizontal);"),v("extra","color: var(--n-item-text-color-active-hover-horizontal);")])])]),R("child-active",[ee(null,[v("icon","color: var(--n-item-icon-color-child-active-hover-horizontal);"),u("menu-item-content-header",`
 color: var(--n-item-text-color-child-active-hover-horizontal);
 `,[g("a","color: var(--n-item-text-color-child-active-hover-horizontal);"),v("extra","color: var(--n-item-text-color-child-active-hover-horizontal);")])])]),ee("border-bottom: 2px solid var(--n-border-color-horizontal);",Fe)]),u("menu-item-content-header",[g("a","color: var(--n-item-text-color-horizontal);")])])]),le("responsive",[u("menu-item-content-header",`
 overflow: hidden;
 text-overflow: ellipsis;
 `)]),R("collapsed",[u("menu-item-content",[R("selected",[g("&::before",`
 background-color: var(--n-item-color-active-collapsed) !important;
 `)]),u("menu-item-content-header","opacity: 0;"),v("arrow","opacity: 0;"),v("icon","color: var(--n-item-icon-color-collapsed);")])]),u("menu-item",`
 height: var(--n-item-height);
 margin-top: 6px;
 position: relative;
 `),u("menu-item-content",`
 box-sizing: border-box;
 line-height: 1.75;
 height: 100%;
 display: grid;
 grid-template-areas: "icon content arrow";
 grid-template-columns: auto 1fr auto;
 align-items: center;
 cursor: pointer;
 position: relative;
 padding-right: 18px;
 transition:
 background-color .3s var(--n-bezier),
 padding-left .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[g("> *","z-index: 1;"),g("&::before",`
 z-index: auto;
 content: "";
 background-color: #0000;
 position: absolute;
 left: 8px;
 right: 8px;
 top: 0;
 bottom: 0;
 pointer-events: none;
 border-radius: var(--n-border-radius);
 transition: background-color .3s var(--n-bezier);
 `),R("disabled",`
 opacity: .45;
 cursor: not-allowed;
 `),R("collapsed",[v("arrow","transform: rotate(0);")]),R("selected",[g("&::before","background-color: var(--n-item-color-active);"),v("arrow","color: var(--n-arrow-color-active);"),v("icon","color: var(--n-item-icon-color-active);"),u("menu-item-content-header",`
 color: var(--n-item-text-color-active);
 `,[g("a","color: var(--n-item-text-color-active);"),v("extra","color: var(--n-item-text-color-active);")])]),R("child-active",[u("menu-item-content-header",`
 color: var(--n-item-text-color-child-active);
 `,[g("a",`
 color: var(--n-item-text-color-child-active);
 `),v("extra",`
 color: var(--n-item-text-color-child-active);
 `)]),v("arrow",`
 color: var(--n-arrow-color-child-active);
 `),v("icon",`
 color: var(--n-item-icon-color-child-active);
 `)]),le("disabled",[le("selected, child-active",[g("&:focus-within",Le)]),R("selected",[ee(null,[v("arrow","color: var(--n-arrow-color-active-hover);"),v("icon","color: var(--n-item-icon-color-active-hover);"),u("menu-item-content-header",`
 color: var(--n-item-text-color-active-hover);
 `,[g("a","color: var(--n-item-text-color-active-hover);"),v("extra","color: var(--n-item-text-color-active-hover);")])])]),R("child-active",[ee(null,[v("arrow","color: var(--n-arrow-color-child-active-hover);"),v("icon","color: var(--n-item-icon-color-child-active-hover);"),u("menu-item-content-header",`
 color: var(--n-item-text-color-child-active-hover);
 `,[g("a","color: var(--n-item-text-color-child-active-hover);"),v("extra","color: var(--n-item-text-color-child-active-hover);")])])]),R("selected",[ee(null,[g("&::before","background-color: var(--n-item-color-active-hover);")])]),ee(null,Le)]),v("icon",`
 grid-area: icon;
 color: var(--n-item-icon-color);
 transition:
 color .3s var(--n-bezier),
 font-size .3s var(--n-bezier),
 margin-right .3s var(--n-bezier);
 box-sizing: content-box;
 display: inline-flex;
 align-items: center;
 justify-content: center;
 `),v("arrow",`
 grid-area: arrow;
 font-size: 16px;
 color: var(--n-arrow-color);
 transform: rotate(180deg);
 opacity: 1;
 transition:
 color .3s var(--n-bezier),
 transform 0.2s var(--n-bezier),
 opacity 0.2s var(--n-bezier);
 `),u("menu-item-content-header",`
 grid-area: content;
 transition:
 color .3s var(--n-bezier),
 opacity .3s var(--n-bezier);
 opacity: 1;
 white-space: nowrap;
 color: var(--n-item-text-color);
 `,[g("a",`
 outline: none;
 text-decoration: none;
 transition: color .3s var(--n-bezier);
 color: var(--n-item-text-color);
 `,[g("&::before",`
 content: "";
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `)]),v("extra",`
 font-size: .93em;
 color: var(--n-group-text-color);
 transition: color .3s var(--n-bezier);
 `)])]),u("submenu",`
 cursor: pointer;
 position: relative;
 margin-top: 6px;
 `,[u("menu-item-content",`
 height: var(--n-item-height);
 `),u("submenu-children",`
 overflow: hidden;
 padding: 0;
 `,[po({duration:".2s"})])]),u("menu-item-group",[u("menu-item-group-title",`
 margin-top: 6px;
 color: var(--n-group-text-color);
 cursor: default;
 font-size: .93em;
 height: 36px;
 display: flex;
 align-items: center;
 transition:
 padding-left .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `)])]),u("menu-tooltip",[g("a",`
 color: inherit;
 text-decoration: none;
 `)]),u("menu-divider",`
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-divider-color);
 height: 1px;
 margin: 6px 18px;
 `)]);function ee(e,t){return[R("hover",e,t),g("&:hover",e,t)]}var lt=L({name:"MenuDivider",setup(){const{mergedClsPrefixRef:e,isHorizontalRef:t}=W(de);return()=>t.value?null:(m(),y("div",{key:1,class:_(`${e.value}-menu-divider`)},null,2))}}),at=L({name:"ChevronDownFilled",render(){return(()=>{const e=Ie("f3af82a2aab086a5");return e[0]||(e[0]=G("svg",{viewBox:"0 0 16 16",fill:"none",xmlns:"http://www.w3.org/2000/svg"},[G("path",{d:"M3.20041 5.73966C3.48226 5.43613 3.95681 5.41856 4.26034 5.70041L8 9.22652L11.7397 5.70041C12.0432 5.41856 12.5177 5.43613 12.7996 5.73966C13.0815 6.0432 13.0639 6.51775 12.7603 6.7996L8.51034 10.7996C8.22258 11.0668 7.77743 11.0668 7.48967 10.7996L3.23966 6.7996C2.93613 6.51775 2.91856 6.0432 3.20041 5.73966Z",fill:"currentColor"})],-1))})()}});const ct=["onClick"];var We=L({name:"MenuOptionContent",props:{collapsed:Boolean,disabled:Boolean,title:[String,Function],icon:Function,extra:[String,Function],showArrow:Boolean,childActive:Boolean,hover:Boolean,paddingLeft:Number,selected:Boolean,maxIconSize:{type:Number,required:!0},activeIconSize:{type:Number,required:!0},iconMarginRight:{type:Number,required:!0},clsPrefix:{type:String,required:!0},onClick:Function,tmNode:{type:Object,required:!0},isEllipsisPlaceholder:Boolean},setup(e){const{props:t}=W(de);return{menuProps:t,style:b(()=>{const{paddingLeft:o}=e;return{paddingLeft:o&&`${o}px`}}),iconStyle:b(()=>{const{maxIconSize:o,activeIconSize:r,iconMarginRight:a}=e;return{width:`${o}px`,height:`${o}px`,fontSize:`${r}px`,marginRight:`${a}px`}})}},render(){const{clsPrefix:e,tmNode:t,menuProps:{renderIcon:o,renderLabel:r,renderExtra:a,expandIcon:i}}=this,s=o?o(t.rawNode):re(this.icon);return(()=>{const d=Ie("7bb10afc6caf8fa4");return m(),y("div",{onClick:c=>{var h;(h=this.onClick)==null||h.call(this,c)},role:"none",class:_([`${e}-menu-item-content`,{[`${e}-menu-item-content--selected`]:this.selected,[`${e}-menu-item-content--collapsed`]:this.collapsed,[`${e}-menu-item-content--child-active`]:this.childActive,[`${e}-menu-item-content--disabled`]:this.disabled,[`${e}-menu-item-content--hover`]:this.hover}]),style:U(this.style)},[I(()=>s&&(m(),y("div",{class:_(`${e}-menu-item-content__icon`),style:U(this.iconStyle),role:"none"},[I(()=>[s])],6))),G("div",{class:_(`${e}-menu-item-content-header`),role:"none"},[this.isEllipsisPlaceholder?(m(),y(X,{key:0},[I(()=>this.title)],64)):(m(),y(X,{key:1},[r?(m(),y(X,{key:0},[I(()=>r(t.rawNode))],64)):(m(),y(X,{key:1},[I(()=>re(this.title))],64))],64)),this.extra||a?(m(),y("span",{key:2,class:_(`${e}-menu-item-content-header__extra`)},[d[0]||(d[0]=I(" ",-1)),a?(m(),y(X,{key:0},[I(()=>a(t.rawNode))],64)):(m(),y(X,{key:1},[I(()=>re(this.extra))],64))],2)):I(()=>null)],2),this.showArrow?(m(),$(De,{key:0,ariaHidden:!0,class:_(`${e}-menu-item-content__arrow`),clsPrefix:e},{default:()=>i?i(t.rawNode):(m(),$(at,{key:1}))},1032,["class","clsPrefix"])):I(()=>null)],14,ct)})()}});const ve=8;function He(e){const t=W(de),{props:o,mergedCollapsedRef:r}=t,a=W(Xe,null),i=W(Re,null),s=b(()=>o.mode==="horizontal"),d=b(()=>s.value?o.dropdownPlacement:"tmNodes"in e?"right-start":"right"),c=b(()=>Math.max(o.collapsedIconSize??o.iconSize,o.iconSize));return{dropdownPlacement:d,activeIconSize:b(()=>!s.value&&e.root&&r.value?o.collapsedIconSize??o.iconSize:o.iconSize),maxIconSize:c,paddingLeft:b(()=>{if(s.value)return;const{collapsedWidth:h,indent:H,rootIndent:z}=o,{root:p,isGroup:S}=e,P=z===void 0?H:z;return p?r.value?h/2-c.value/2:P:i&&typeof i.paddingLeftRef.value=="number"?r.value?h/2-c.value/2:H/2+i.paddingLeftRef.value:a&&typeof a.paddingLeftRef.value=="number"?(S?H/2:H)+a.paddingLeftRef.value:0}),iconMarginRight:b(()=>{const{collapsedWidth:h,indent:H,rootIndent:z}=o,{value:p}=c,{root:S}=e;return s.value||!S||!r.value?ve:(z===void 0?H:z)+p+ve-(h+p)/2}),NMenu:t,NSubmenu:a,NMenuOptionGroup:i}}const ke={internalKey:{type:[String,Number],required:!0},root:Boolean,isGroup:Boolean,level:{type:Number,required:!0},title:[String,Function],extra:[String,Function]},Je={...ke,tmNode:{type:Object,required:!0},disabled:Boolean,icon:Function,onClick:Function},st=Se(Je),dt=L({name:"MenuOption",props:Je,setup(e){const t=He(e),{NSubmenu:o,NMenu:r,NMenuOptionGroup:a}=t,{props:i,mergedClsPrefixRef:s,mergedCollapsedRef:d}=r,c=o?o.mergedDisabledRef:a?a.mergedDisabledRef:{value:!1},h=b(()=>c.value||e.disabled);function H(p){const{onClick:S}=e;S&&S(p)}function z(p){h.value||(r.doSelect(e.internalKey,e.tmNode.rawNode),H(p))}return{mergedClsPrefix:s,dropdownPlacement:t.dropdownPlacement,paddingLeft:t.paddingLeft,iconMarginRight:t.iconMarginRight,maxIconSize:t.maxIconSize,activeIconSize:t.activeIconSize,mergedTheme:r.mergedThemeRef,menuProps:i,dropdownEnabled:be(()=>e.root&&d.value&&i.mode!=="horizontal"&&!h.value),selected:be(()=>r.mergedValueRef.value===e.internalKey),mergedDisabled:h,handleClick:z}},render(){const{mergedClsPrefix:e,mergedTheme:t,tmNode:o,menuProps:{renderLabel:r,nodeProps:a}}=this,i=a==null?void 0:a(o.rawNode);return m(),y("div",Z(i,{role:"menuitem",class:[`${e}-menu-item`,i==null?void 0:i.class]}),[(m(),$(Ho,{theme:t.peers.Tooltip,themeOverrides:t.peerOverrides.Tooltip,trigger:"hover",placement:this.dropdownPlacement,disabled:!this.dropdownEnabled||this.title===void 0,internalExtraClass:["menu-tooltip"]},{default:()=>r?r(o.rawNode):re(this.title),trigger:()=>(m(),$(We,{tmNode:o,clsPrefix:e,paddingLeft:this.paddingLeft,iconMarginRight:this.iconMarginRight,maxIconSize:this.maxIconSize,activeIconSize:this.activeIconSize,selected:this.selected,title:this.title,extra:this.extra,disabled:this.mergedDisabled,icon:this.icon,onClick:this.handleClick},null,8,["tmNode","clsPrefix","paddingLeft","iconMarginRight","maxIconSize","activeIconSize","selected","title","extra","disabled","icon","onClick"]))},1032,["theme","themeOverrides","placement","disabled"]))],16)}}),Qe={...ke,tmNode:{type:Object,required:!0},tmNodes:{type:Array,required:!0}},ut=Se(Qe),vt=L({name:"MenuOptionGroup",props:Qe,setup(e){const t=He(e),{NSubmenu:o}=t,r=b(()=>o!=null&&o.mergedDisabledRef.value?!0:e.tmNode.disabled);oe(Re,{paddingLeftRef:t.paddingLeft,mergedDisabledRef:r});const{mergedClsPrefixRef:a,props:i}=W(de);return function(){const{value:s}=a,d=t.paddingLeft.value,{nodeProps:c}=i,h=c==null?void 0:c(e.tmNode.rawNode);return(()=>{const H=Ie("45eca6a63be5028b");return m(),y("div",{class:_(`${s}-menu-item-group`),role:"group"},[G("div",Z(h,{class:[`${s}-menu-item-group-title`,h==null?void 0:h.class],style:[(h==null?void 0:h.style)||"",d!==void 0?`padding-left: ${d}px;`:""]}),[I(()=>re(e.title)),e.extra?(m(),y(X,{key:0},[H[0]||(H[0]=I(" ",-1)),I(()=>re(e.extra))],64)):I(()=>null)],16),G("div",null,[I(()=>e.tmNodes.map(z=>_e(z,i)))])],2)})()}}}),mt=["aria-expanded","id"],ht=["aria-expanded","id"],Ze={...ke,rawNodes:{type:Array,default:()=>[]},tmNodes:{type:Array,default:()=>[]},tmNode:{type:Object,required:!0},disabled:Boolean,icon:Function,onClick:Function,domId:String,virtualChildActive:{type:Boolean,default:void 0},isEllipsisPlaceholder:Boolean},ft=Se(Ze),xe=L({name:"Submenu",props:Ze,setup(e){const t=He(e),{NMenu:o,NSubmenu:r}=t,{props:a,mergedCollapsedRef:i,mergedThemeRef:s}=o,d=b(()=>{const{disabled:p}=e;return r!=null&&r.mergedDisabledRef.value||a.disabled?!0:p}),c=O(!1);oe(Xe,{paddingLeftRef:t.paddingLeft,mergedDisabledRef:d}),oe(Re,null);function h(){const{onClick:p}=e;p&&p()}function H(){d.value||(i.value||o.toggleExpand(e.internalKey),h())}function z(p){c.value=p}return{menuProps:a,mergedTheme:s,doSelect:o.doSelect,inverted:o.invertedRef,isHorizontal:o.isHorizontalRef,mergedClsPrefix:o.mergedClsPrefixRef,maxIconSize:t.maxIconSize,activeIconSize:t.activeIconSize,iconMarginRight:t.iconMarginRight,dropdownPlacement:t.dropdownPlacement,dropdownShow:c,paddingLeft:t.paddingLeft,mergedDisabled:d,mergedValue:o.mergedValueRef,childActive:be(()=>e.virtualChildActive??o.activePathRef.value.includes(e.internalKey)),collapsed:b(()=>a.mode==="horizontal"?!1:i.value?!0:!o.mergedExpandedKeysRef.value.includes(e.internalKey)),dropdownEnabled:b(()=>!d.value&&(a.mode==="horizontal"||i.value)),handlePopoverShowChange:z,handleClick:H}},render(){var i;const{mergedClsPrefix:e,menuProps:{renderIcon:t,renderLabel:o}}=this,r=()=>{const{isHorizontal:s,paddingLeft:d,collapsed:c,mergedDisabled:h,maxIconSize:H,activeIconSize:z,title:p,childActive:S,icon:P,handleClick:x,menuProps:{nodeProps:T},dropdownShow:B,iconMarginRight:ie,tmNode:J,mergedClsPrefix:K,isEllipsisPlaceholder:N,extra:A}=this,w=T==null?void 0:T(J.rawNode);return m(),y("div",Z(w,{class:[`${K}-menu-item`,w==null?void 0:w.class],role:"menuitem"}),[(m(),$(We,{tmNode:J,paddingLeft:d,collapsed:c,disabled:h,iconMarginRight:ie,maxIconSize:H,activeIconSize:z,title:p,extra:A,showArrow:!s,childActive:S,clsPrefix:K,icon:P,hover:B,onClick:x,isEllipsisPlaceholder:N},null,8,["tmNode","paddingLeft","collapsed","disabled","iconMarginRight","maxIconSize","activeIconSize","title","extra","showArrow","childActive","clsPrefix","icon","hover","onClick","isEllipsisPlaceholder"]))],16)},a=()=>(m(),$(go,null,{default:()=>{const{tmNodes:s,collapsed:d}=this;return d?null:(m(),y("div",{key:1,class:_(`${e}-submenu-children`),role:"menu"},[I(()=>s.map(c=>_e(c,this.menuProps)))],2))}},1024));return this.root?(m(),$(Ue,Z({key:2,size:"large",trigger:"hover"},(i=this.menuProps)==null?void 0:i.dropdownProps,{themeOverrides:this.mergedTheme.peerOverrides.Dropdown,theme:this.mergedTheme.peers.Dropdown,builtinThemeOverrides:{fontSizeLarge:"14px",optionIconSizeLarge:"18px"},value:this.mergedValue,disabled:!this.dropdownEnabled,placement:this.dropdownPlacement,keyField:this.menuProps.keyField,labelField:this.menuProps.labelField,childrenField:this.menuProps.childrenField,onUpdateShow:this.handlePopoverShowChange,options:this.rawNodes,onSelect:this.doSelect,inverted:this.inverted,renderIcon:t,renderLabel:o}),{default:()=>(m(),y("div",{class:_(`${e}-submenu`),role:"menu","aria-expanded":!this.collapsed,id:this.domId},[I(()=>r()),this.isHorizontal?I(()=>null):(m(),y(X,{key:1},[I(()=>a())],64))],10,mt))},1040,["themeOverrides","theme","value","disabled","placement","keyField","labelField","childrenField","onUpdateShow","options","onSelect","inverted","renderIcon","renderLabel"])):(m(),y("div",{key:3,class:_(`${e}-submenu`),role:"menu","aria-expanded":!this.collapsed,id:this.domId},[I(()=>r()),I(()=>a())],10,ht))}});function ye(e){return e.type==="divider"||e.type==="render"}function pt(e){return e.type==="divider"}function _e(e,t){const{rawNode:o}=e,{show:r}=o;if(r===!1)return null;if(ye(o))return pt(o)?(m(),$(lt,Z({key:e.key},o.props),null,16)):null;const{labelField:a}=t,{key:i,level:s,isGroup:d}=e,c={...o,title:o.title||o[a],extra:o.titleExtra||o.extra,key:i,internalKey:i,level:s,root:s===0,isGroup:d};return e.children?e.isGroup?te(vt,fe(c,ut,{tmNode:e,tmNodes:e.children,key:i})):te(xe,fe(c,ft,{key:i,rawNodes:o[t.childrenField],tmNodes:e.children,tmNode:e})):te(dt,fe(c,st,{key:i,tmNode:e}))}const gt={...j.props,options:{type:Array,default:()=>[]},collapsed:{type:Boolean,default:void 0},collapsedWidth:{type:Number,default:48},iconSize:{type:Number,default:20},collapsedIconSize:{type:Number,default:24},rootIndent:Number,indent:{type:Number,default:32},labelField:{type:String,default:"label"},keyField:{type:String,default:"key"},childrenField:{type:String,default:"children"},disabledField:{type:String,default:"disabled"},defaultExpandAll:Boolean,defaultExpandedKeys:Array,expandedKeys:Array,value:[String,Number],defaultValue:{type:[String,Number],default:null},mode:{type:String,default:"vertical"},watchProps:{type:Array,default:void 0},disabled:Boolean,show:{type:Boolean,default:!0},inverted:Boolean,"onUpdate:expandedKeys":[Function,Array],onUpdateExpandedKeys:[Function,Array],onUpdateValue:[Function,Array],"onUpdate:value":[Function,Array],expandIcon:Function,renderIcon:Function,renderLabel:Function,renderExtra:Function,dropdownProps:Object,accordion:Boolean,nodeProps:Function,dropdownPlacement:{type:String,default:"bottom"},responsive:Boolean,items:Array,onOpenNamesChange:[Function,Array],onSelect:[Function,Array],onExpandedNamesChange:[Function,Array],expandedNames:Array,defaultExpandedNames:Array};var bt=L({name:"Menu",inheritAttrs:!1,props:gt,setup(e){const{mergedClsPrefixRef:t,inlineThemeDisabled:o}=ce(e),r=j("Menu","-menu",it,Do,e,t),a=W(Ge,null),i=b(()=>{const{collapsed:f}=e;if(f!==void 0)return f;if(a){const{collapseModeRef:k,collapsedRef:n}=a;if(k.value==="width")return n.value??!1}return!1}),s=b(()=>{const{keyField:f,childrenField:k,disabledField:n}=e;return pe(e.items||e.options,{getIgnored(C){return ye(C)},getChildren(C){return C[k]},getDisabled(C){return C[n]},getKey(C){return C[f]??C.name}})}),d=b(()=>new Set(s.value.treeNodes.map(f=>f.key))),{watchProps:c}=e,h=O(null);c!=null&&c.includes("defaultValue")?Ne(()=>{h.value=e.defaultValue}):h.value=e.defaultValue;const H=ae(e,"value"),z=Ce(H,h),p=O([]),S=()=>{p.value=e.defaultExpandAll?s.value.getNonLeafKeys():e.defaultExpandedNames||e.defaultExpandedKeys||s.value.getPath(z.value,{includeSelf:!1}).keyPath};c!=null&&c.includes("defaultExpandedKeys")?Ne(S):S();const P=_o(e,["expandedNames","expandedKeys"]),x=Ce(P,p),T=b(()=>s.value.treeNodes),B=b(()=>s.value.getPath(z.value).keyPath);oe(de,{props:e,mergedCollapsedRef:i,mergedThemeRef:r,mergedValueRef:z,mergedExpandedKeysRef:x,activePathRef:B,mergedClsPrefixRef:t,isHorizontalRef:b(()=>e.mode==="horizontal"),invertedRef:ae(e,"inverted"),doSelect:ie,toggleExpand:K});function ie(f,k){const{"onUpdate:value":n,onUpdateValue:C,onSelect:E}=e;C&&V(C,f,k),n&&V(n,f,k),E&&V(E,f,k),h.value=f}function J(f){const{"onUpdate:expandedKeys":k,onUpdateExpandedKeys:n,onExpandedNamesChange:C,onOpenNamesChange:E}=e;k&&V(k,f),n&&V(n,f),C&&V(C,f),E&&V(E,f),p.value=f}function K(f){const k=Array.from(x.value),n=k.findIndex(C=>C===f);if(~n)k.splice(n,1);else{if(e.accordion&&d.value.has(f)){const C=k.findIndex(E=>d.value.has(E));C>-1&&k.splice(C,1)}k.push(f)}J(k)}const N=f=>{const k=s.value.getPath(f??z.value,{includeSelf:!1}).keyPath;if(!k.length)return;const n=Array.from(x.value),C=new Set([...n,...k]);e.accordion&&d.value.forEach(E=>{C.has(E)&&!k.includes(E)&&C.delete(E)}),J(Array.from(C))},A=b(()=>{const{inverted:f}=e,{common:{cubicBezierEaseInOut:k},self:n}=r.value,{borderRadius:C,borderColorHorizontal:E,fontSize:ao,itemHeight:co,dividerColor:so}=n,l={"--n-divider-color":so,"--n-bezier":k,"--n-font-size":ao,"--n-border-color-horizontal":E,"--n-border-radius":C,"--n-item-height":co};return f?(l["--n-group-text-color"]=n.groupTextColorInverted,l["--n-color"]=n.colorInverted,l["--n-item-text-color"]=n.itemTextColorInverted,l["--n-item-text-color-hover"]=n.itemTextColorHoverInverted,l["--n-item-text-color-active"]=n.itemTextColorActiveInverted,l["--n-item-text-color-child-active"]=n.itemTextColorChildActiveInverted,l["--n-item-text-color-child-active-hover"]=n.itemTextColorChildActiveInverted,l["--n-item-text-color-active-hover"]=n.itemTextColorActiveHoverInverted,l["--n-item-icon-color"]=n.itemIconColorInverted,l["--n-item-icon-color-hover"]=n.itemIconColorHoverInverted,l["--n-item-icon-color-active"]=n.itemIconColorActiveInverted,l["--n-item-icon-color-active-hover"]=n.itemIconColorActiveHoverInverted,l["--n-item-icon-color-child-active"]=n.itemIconColorChildActiveInverted,l["--n-item-icon-color-child-active-hover"]=n.itemIconColorChildActiveHoverInverted,l["--n-item-icon-color-collapsed"]=n.itemIconColorCollapsedInverted,l["--n-item-text-color-horizontal"]=n.itemTextColorHorizontalInverted,l["--n-item-text-color-hover-horizontal"]=n.itemTextColorHoverHorizontalInverted,l["--n-item-text-color-active-horizontal"]=n.itemTextColorActiveHorizontalInverted,l["--n-item-text-color-child-active-horizontal"]=n.itemTextColorChildActiveHorizontalInverted,l["--n-item-text-color-child-active-hover-horizontal"]=n.itemTextColorChildActiveHoverHorizontalInverted,l["--n-item-text-color-active-hover-horizontal"]=n.itemTextColorActiveHoverHorizontalInverted,l["--n-item-icon-color-horizontal"]=n.itemIconColorHorizontalInverted,l["--n-item-icon-color-hover-horizontal"]=n.itemIconColorHoverHorizontalInverted,l["--n-item-icon-color-active-horizontal"]=n.itemIconColorActiveHorizontalInverted,l["--n-item-icon-color-active-hover-horizontal"]=n.itemIconColorActiveHoverHorizontalInverted,l["--n-item-icon-color-child-active-horizontal"]=n.itemIconColorChildActiveHorizontalInverted,l["--n-item-icon-color-child-active-hover-horizontal"]=n.itemIconColorChildActiveHoverHorizontalInverted,l["--n-arrow-color"]=n.arrowColorInverted,l["--n-arrow-color-hover"]=n.arrowColorHoverInverted,l["--n-arrow-color-active"]=n.arrowColorActiveInverted,l["--n-arrow-color-active-hover"]=n.arrowColorActiveHoverInverted,l["--n-arrow-color-child-active"]=n.arrowColorChildActiveInverted,l["--n-arrow-color-child-active-hover"]=n.arrowColorChildActiveHoverInverted,l["--n-item-color-hover"]=n.itemColorHoverInverted,l["--n-item-color-active"]=n.itemColorActiveInverted,l["--n-item-color-active-hover"]=n.itemColorActiveHoverInverted,l["--n-item-color-active-collapsed"]=n.itemColorActiveCollapsedInverted):(l["--n-group-text-color"]=n.groupTextColor,l["--n-color"]=n.color,l["--n-item-text-color"]=n.itemTextColor,l["--n-item-text-color-hover"]=n.itemTextColorHover,l["--n-item-text-color-active"]=n.itemTextColorActive,l["--n-item-text-color-child-active"]=n.itemTextColorChildActive,l["--n-item-text-color-child-active-hover"]=n.itemTextColorChildActiveHover,l["--n-item-text-color-active-hover"]=n.itemTextColorActiveHover,l["--n-item-icon-color"]=n.itemIconColor,l["--n-item-icon-color-hover"]=n.itemIconColorHover,l["--n-item-icon-color-active"]=n.itemIconColorActive,l["--n-item-icon-color-active-hover"]=n.itemIconColorActiveHover,l["--n-item-icon-color-child-active"]=n.itemIconColorChildActive,l["--n-item-icon-color-child-active-hover"]=n.itemIconColorChildActiveHover,l["--n-item-icon-color-collapsed"]=n.itemIconColorCollapsed,l["--n-item-text-color-horizontal"]=n.itemTextColorHorizontal,l["--n-item-text-color-hover-horizontal"]=n.itemTextColorHoverHorizontal,l["--n-item-text-color-active-horizontal"]=n.itemTextColorActiveHorizontal,l["--n-item-text-color-child-active-horizontal"]=n.itemTextColorChildActiveHorizontal,l["--n-item-text-color-child-active-hover-horizontal"]=n.itemTextColorChildActiveHoverHorizontal,l["--n-item-text-color-active-hover-horizontal"]=n.itemTextColorActiveHoverHorizontal,l["--n-item-icon-color-horizontal"]=n.itemIconColorHorizontal,l["--n-item-icon-color-hover-horizontal"]=n.itemIconColorHoverHorizontal,l["--n-item-icon-color-active-horizontal"]=n.itemIconColorActiveHorizontal,l["--n-item-icon-color-active-hover-horizontal"]=n.itemIconColorActiveHoverHorizontal,l["--n-item-icon-color-child-active-horizontal"]=n.itemIconColorChildActiveHorizontal,l["--n-item-icon-color-child-active-hover-horizontal"]=n.itemIconColorChildActiveHoverHorizontal,l["--n-arrow-color"]=n.arrowColor,l["--n-arrow-color-hover"]=n.arrowColorHover,l["--n-arrow-color-active"]=n.arrowColorActive,l["--n-arrow-color-active-hover"]=n.arrowColorActiveHover,l["--n-arrow-color-child-active"]=n.arrowColorChildActive,l["--n-arrow-color-child-active-hover"]=n.arrowColorChildActiveHover,l["--n-item-color-hover"]=n.itemColorHover,l["--n-item-color-active"]=n.itemColorActive,l["--n-item-color-active-hover"]=n.itemColorActiveHover,l["--n-item-color-active-collapsed"]=n.itemColorActiveCollapsed),l}),w=o?se("menu",b(()=>e.inverted?"a":"b"),A,e):void 0,Q=Co(),q=O(null),me=O(null);let F=!0;const Te=()=>{var f;F?F=!1:(f=q.value)==null||f.sync({showAllItemsBeforeCalculate:!0})};function eo(){return document.getElementById(Q)}const ue=O(-1);function oo(f){ue.value=e.options.length-f}function to(f){f||(ue.value=-1)}const ro=b(()=>{const f=ue.value;return{children:f===-1?[]:e.options.slice(f)}}),no=b(()=>{const{childrenField:f,disabledField:k,keyField:n}=e;return pe([ro.value],{getIgnored(C){return ye(C)},getChildren(C){return C[f]},getDisabled(C){return C[k]},getKey(C){return C[n]??C.name}})}),io=b(()=>pe([{}]).treeNodes[0]);function lo(){var C;if(ue.value===-1)return m(),$(xe,{root:!0,level:0,key:"__ellpisisGroupPlaceholder__",internalKey:"__ellpisisGroupPlaceholder__",title:"···",tmNode:io.value,domId:Q,isEllipsisPlaceholder:!0},null,8,["tmNode","domId"]);const f=no.value.treeNodes[0],k=B.value,n=!!((C=f.children)!=null&&C.some(E=>k.includes(E.key)));return m(),$(xe,{level:0,root:!0,key:"__ellpisisGroup__",internalKey:"__ellpisisGroup__",title:"···",virtualChildActive:n,tmNode:f,domId:Q,rawNodes:f.rawNode.children||[],tmNodes:f.children||[],isEllipsisPlaceholder:!0},null,8,["virtualChildActive","tmNode","domId","rawNodes","tmNodes"])}return{mergedClsPrefix:t,controlledExpandedKeys:P,uncontrolledExpanededKeys:p,mergedExpandedKeys:x,uncontrolledValue:h,mergedValue:z,activePath:B,tmNodes:T,mergedTheme:r,mergedCollapsed:i,cssVars:o?void 0:A,themeClass:w==null?void 0:w.themeClass,overflowRef:q,counterRef:me,updateCounter:()=>{},onResize:Te,onUpdateOverflow:to,onUpdateCount:oo,renderCounter:lo,getCounter:eo,onRender:w==null?void 0:w.onRender,showOption:N,deriveResponsiveState:Te}},render(){const{mergedClsPrefix:e,mode:t,themeClass:o,onRender:r}=this;r==null||r();const a=()=>this.tmNodes.map(d=>_e(d,this.$props)),i=t==="horizontal"&&this.responsive,s=()=>te("div",Z(this.$attrs,{role:t==="horizontal"?"menubar":"menu",class:[`${e}-menu`,o,`${e}-menu--${t}`,i&&`${e}-menu--responsive`,this.mergedCollapsed&&`${e}-menu--collapsed`],style:this.cssVars}),i?(m(),$(ko,{key:2,ref:"overflowRef",onUpdateOverflow:this.onUpdateOverflow,getCounter:this.getCounter,onUpdateCount:this.onUpdateCount,updateCounter:this.updateCounter,style:{width:"100%",display:"flex",overflow:"hidden"}},{default:a,counter:this.renderCounter},1032,["onUpdateOverflow","getCounter","onUpdateCount","updateCounter"])):a());return i?(m(),$(bo,{key:3,onResize:this.onResize},{default:s},1032,["onResize"])):s()}});const Ct={class:"logo"},xt={key:0},yt={key:1},zt={class:"header-left"},It=L({__name:"AdminLayout",setup(e){const t=zo(),o=So(),r=xo(),a=O(!1),i={HomeFilled:"🏠",Setting:"⚙️",User:"👤",Avatar:"🎭",Lock:"🔒",Menu:"📋"},s=p=>()=>te("span",{style:"margin-right:6px"},p&&i[p]||"📄");function d(p){return(p??[]).map(S=>{var P;return{label:S.name,key:S.path,icon:s(S.icon),children:(P=S.children)!=null&&P.length?d(S.children):void 0}})}const c=b(()=>d(r.menus)),h=b(()=>t.path),H=[{label:"退出登录",key:"logout"}];async function z(p){p==="logout"&&(await r.logout(),o.push("/login"))}return(p,S)=>{const P=Io("router-view");return m(),$(M(Ee),{"has-sider":"",style:{height:"100vh"}},{default:Y(()=>[D(M(nt),{bordered:"","collapse-mode":"width","collapsed-width":64,width:220,collapsed:a.value,"show-trigger":"",onCollapse:S[0]||(S[0]=x=>a.value=!0),onExpand:S[1]||(S[1]=x=>a.value=!1)},{default:Y(()=>[G("div",Ct,[a.value?(m(),y("span",yt,"S")):(m(),y("span",xt,"SmileX Admin"))]),D(M(bt),{collapsed:a.value,"collapsed-width":64,options:c.value,value:h.value},null,8,["collapsed","options","value"])]),_:1},8,["collapsed"]),D(M(Ee),null,{default:Y(()=>[D(M(Wo),{bordered:"",class:"header"},{default:Y(()=>[G("div",zt,[D(M(Lo),null,{default:Y(()=>[D(M(Oo),null,{default:Y(()=>{var x;return[Be($e(((x=M(t).meta)==null?void 0:x.title)||"首页"),1)]}),_:1})]),_:1})]),D(M(Ue),{options:H,onSelect:z},{default:Y(()=>[D(M(yo),{quaternary:""},{default:Y(()=>{var x,T;return[Be($e(((x=M(r).user)==null?void 0:x.nickname)||((T=M(r).user)==null?void 0:T.username)),1)]}),_:1})]),_:1})]),_:1}),D(M(qo),{"content-style":"padding: 16px;","native-scrollbar":!1},{default:Y(()=>[D(P)]),_:1})]),_:1})]),_:1})}}}),kt=To(It,[["__scopeId","data-v-7454bf7f"]]);export{kt as default};
