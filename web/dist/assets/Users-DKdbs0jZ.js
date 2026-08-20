import{d as $e,T as Re,A as ce,B as o,U as he,W as ne,X as $,Y as fe,b as me,C as ge,Z as oe,e as p,f as D,F as K,n as B,D as O,E as w,l as j,g as ve,u as Ve,_ as Fe,r as v,I as Ue,c as X,$ as se,a0 as Pe,m as We,a1 as Me,a2 as T,p as re,a3 as R,z as Oe,s as Te,a4 as De,N as r,O as u,Q as n,a5 as Ke,a6 as V,a7 as L,a8 as Le,R as Ne,a9 as ue,aa as pe,G as Ie,ab as de,K as je,L as He,ac as Ae,ad as Ee,ae as Xe,af as Ye,k as A,ag as qe}from"./index-C-2PX-C-.js";import{D as Ge,S as be,a as Qe,T as Ze}from"./Space-CIqFIzT1.js";import{u as Je,I as N,F as ea,a as I}from"./FormItem-T-2ecOeW.js";import{u as aa}from"./format-length-icqZkaSq.js";import"./Dropdown-CM88cknF.js";import"./next-frame-once-C5Ksf8W7.js";import"./get-slot-VaplJ9hC.js";var ta={buttonHeightSmall:"14px",buttonHeightMedium:"18px",buttonHeightLarge:"22px",buttonWidthSmall:"14px",buttonWidthMedium:"18px",buttonWidthLarge:"22px",buttonWidthPressedSmall:"20px",buttonWidthPressedMedium:"24px",buttonWidthPressedLarge:"28px",railHeightSmall:"18px",railHeightMedium:"22px",railHeightLarge:"26px",railWidthSmall:"32px",railWidthMedium:"40px",railWidthLarge:"48px"};function la(t){const{primaryColor:h,opacityDisabled:y,borderRadius:d,textColor3:k}=t;return{...ta,iconColor:k,textColor:"white",loadingColor:h,opacityDisabled:y,railColor:"rgba(0, 0, 0, .14)",railColorActive:h,buttonBoxShadow:"0 1px 4px 0 rgba(0, 0, 0, 0.3), inset 0 0 1px 0 rgba(0, 0, 0, 0.05)",buttonColor:"#FFF",railBorderRadiusSmall:d,railBorderRadiusMedium:d,railBorderRadiusLarge:d,buttonBorderRadiusSmall:d,buttonBorderRadiusMedium:d,buttonBorderRadiusLarge:d,boxShadowFocus:`0 0 0 2px ${Re(h,{alpha:.2})}`}}const ia={common:$e,self:la};var na=ce("switch",`
 height: var(--n-height);
 min-width: var(--n-width);
 vertical-align: middle;
 user-select: none;
 -webkit-user-select: none;
 display: inline-flex;
 outline: none;
 justify-content: center;
 align-items: center;
`,[o("children-placeholder",`
 height: var(--n-rail-height);
 display: flex;
 flex-direction: column;
 overflow: hidden;
 pointer-events: none;
 visibility: hidden;
 `),o("rail-placeholder",`
 display: flex;
 flex-wrap: none;
 `),o("button-placeholder",`
 width: calc(1.75 * var(--n-rail-height));
 height: var(--n-rail-height);
 `),ce("base-loading",`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 font-size: calc(var(--n-button-width) - 4px);
 color: var(--n-loading-color);
 transition: color .3s var(--n-bezier);
 `,[he({left:"50%",top:"50%",originalTransform:"translateX(-50%) translateY(-50%)"})]),o("checked, unchecked",`
 transition: color .3s var(--n-bezier);
 color: var(--n-text-color);
 box-sizing: border-box;
 position: absolute;
 white-space: nowrap;
 top: 0;
 bottom: 0;
 display: flex;
 align-items: center;
 line-height: 1;
 `),o("checked",`
 right: 0;
 padding-right: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),o("unchecked",`
 left: 0;
 justify-content: flex-end;
 padding-left: calc(1.25 * var(--n-rail-height) - var(--n-offset));
 `),ne("&:focus",[o("rail",`
 box-shadow: var(--n-box-shadow-focus);
 `)]),$("round",[o("rail","border-radius: calc(var(--n-rail-height) / 2);",[o("button","border-radius: calc(var(--n-button-height) / 2);")])]),fe("disabled",[fe("icon",[$("rubber-band",[$("pressed",[o("rail",[o("button","max-width: var(--n-button-width-pressed);")])]),o("rail",[ne("&:active",[o("button","max-width: var(--n-button-width-pressed);")])]),$("active",[$("pressed",[o("rail",[o("button","left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));")])]),o("rail",[ne("&:active",[o("button","left: calc(100% - var(--n-offset) - var(--n-button-width-pressed));")])])])])])]),$("active",[o("rail",[o("button","left: calc(100% - var(--n-button-width) - var(--n-offset))")])]),o("rail",`
 overflow: hidden;
 height: var(--n-rail-height);
 min-width: var(--n-rail-width);
 border-radius: var(--n-rail-border-radius);
 cursor: pointer;
 position: relative;
 transition:
 opacity .3s var(--n-bezier),
 background .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 background-color: var(--n-rail-color);
 `,[o("button-icon",`
 color: var(--n-icon-color);
 transition: color .3s var(--n-bezier);
 font-size: calc(var(--n-button-height) - 4px);
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 display: flex;
 justify-content: center;
 align-items: center;
 line-height: 1;
 `,[he()]),o("button",`
 align-items: center; 
 top: var(--n-offset);
 left: var(--n-offset);
 height: var(--n-button-height);
 width: var(--n-button-width-pressed);
 max-width: var(--n-button-width);
 border-radius: var(--n-button-border-radius);
 background-color: var(--n-button-color);
 box-shadow: var(--n-button-box-shadow);
 box-sizing: border-box;
 cursor: inherit;
 content: "";
 position: absolute;
 transition:
 background-color .3s var(--n-bezier),
 left .3s var(--n-bezier),
 opacity .3s var(--n-bezier),
 max-width .3s var(--n-bezier),
 box-shadow .3s var(--n-bezier);
 `)]),$("active",[o("rail","background-color: var(--n-rail-color-active);")]),$("loading",[o("rail",`
 cursor: wait;
 `)]),$("disabled",[o("rail",`
 cursor: not-allowed;
 opacity: .5;
 `)])]);const oa=["aria-checked","tabindex","onClick","onFocus","onBlur","onKeyup","onKeydown"],sa={...ge.props,size:String,value:{type:[String,Number,Boolean],default:void 0},loading:Boolean,defaultValue:{type:[String,Number,Boolean],default:!1},disabled:{type:Boolean,default:void 0},round:{type:Boolean,default:!0},"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],checkedValue:{type:[String,Number,Boolean],default:!0},uncheckedValue:{type:[String,Number,Boolean],default:!1},railStyle:Function,rubberBand:{type:Boolean,default:!0},spinProps:Object,onChange:[Function,Array]};let E;var ra=me({name:"Switch",props:sa,slots:Object,setup(t){E===void 0&&(typeof CSS<"u"?typeof CSS.supports<"u"?E=CSS.supports("width","max(1px)"):E=!1:E=!0);const{mergedClsPrefixRef:h,inlineThemeDisabled:y,mergedComponentPropsRef:d}=Ve(t),k=ge("Switch","-switch",na,ia,t,h),x=Fe(t,{mergedSize(e){var W,M;if(t.size!==void 0)return t.size;if(e)return e.mergedSize.value;const P=(M=(W=d==null?void 0:d.value)==null?void 0:W.Switch)==null?void 0:M.size;return P||"medium"}}),{mergedSizeRef:f,mergedDisabledRef:S}=x,b=v(t.defaultValue),C=Oe(t,"value"),c=aa(C,b),_=X(()=>c.value===t.checkedValue),s=v(!1),l=v(!1),z=X(()=>{const{railStyle:e}=t;if(e)return e({focused:l.value,checked:_.value})});function m(e){const{"onUpdate:value":P,onChange:W,onUpdateValue:M}=t,{nTriggerFormInput:ee,nTriggerFormChange:ae}=x;P&&se(P,e),M&&se(M,e),W&&se(W,e),b.value=e,ee(),ae()}function Y(){const{nTriggerFormFocus:e}=x;e()}function q(){const{nTriggerFormBlur:e}=x;e()}function G(){t.loading||S.value||(c.value!==t.checkedValue?m(t.checkedValue):m(t.uncheckedValue))}function Q(){l.value=!0,Y()}function Z(){l.value=!1,q(),s.value=!1}function J(e){t.loading||S.value||e.key===" "&&(c.value!==t.checkedValue?m(t.checkedValue):m(t.uncheckedValue),s.value=!1)}function i(e){t.loading||S.value||e.key===" "&&(e.preventDefault(),s.value=!0)}const a=X(()=>{const{value:e}=f,{self:{opacityDisabled:P,railColor:W,railColorActive:M,buttonBoxShadow:ee,buttonColor:ae,boxShadowFocus:we,loadingColor:ye,textColor:ke,iconColor:xe,[T("buttonHeight",e)]:F,[T("buttonWidth",e)]:Se,[T("buttonWidthPressed",e)]:Ce,[T("railHeight",e)]:U,[T("railWidth",e)]:H,[T("railBorderRadius",e)]:_e,[T("buttonBorderRadius",e)]:ze},common:{cubicBezierEaseInOut:Be}}=k.value;let te,le,ie;return E?(te=`calc((${U} - ${F}) / 2)`,le=`max(${U}, ${F})`,ie=`max(${H}, calc(${H} + ${F} - ${U}))`):(te=re((R(U)-R(F))/2),le=re(Math.max(R(U),R(F))),ie=R(U)>R(F)?H:re(R(H)+R(F)-R(U))),{"--n-bezier":Be,"--n-button-border-radius":ze,"--n-button-box-shadow":ee,"--n-button-color":ae,"--n-button-width":Se,"--n-button-width-pressed":Ce,"--n-button-height":F,"--n-height":le,"--n-offset":te,"--n-opacity-disabled":P,"--n-rail-border-radius":_e,"--n-rail-color":W,"--n-rail-color-active":M,"--n-rail-height":U,"--n-rail-width":H,"--n-width":ie,"--n-box-shadow-focus":we,"--n-loading-color":ye,"--n-text-color":ke,"--n-icon-color":xe}}),g=y?Ue("switch",X(()=>f.value[0]),a,t):void 0;return{handleClick:G,handleBlur:Z,handleFocus:Q,handleKeyup:J,handleKeydown:i,mergedRailStyle:z,pressed:s,mergedClsPrefix:h,mergedValue:c,checked:_,mergedDisabled:S,cssVars:y?void 0:a,themeClass:g==null?void 0:g.themeClass,onRender:g==null?void 0:g.onRender}},render(){const{mergedClsPrefix:t,mergedDisabled:h,checked:y,mergedRailStyle:d,onRender:k,$slots:x}=this;k==null||k();const{checked:f,unchecked:S,icon:b,"checked-icon":C,"unchecked-icon":c}=x,_=!(oe(b)&&oe(C)&&oe(c));return p(),D("div",{role:"switch","aria-checked":y,class:w([`${t}-switch`,this.themeClass,_&&`${t}-switch--icon`,y&&`${t}-switch--active`,h&&`${t}-switch--disabled`,this.round&&`${t}-switch--round`,this.loading&&`${t}-switch--loading`,this.pressed&&`${t}-switch--pressed`,this.rubberBand&&`${t}-switch--rubber-band`]),tabindex:this.mergedDisabled?void 0:0,style:ve(this.cssVars),onClick:this.handleClick,onFocus:this.handleFocus,onBlur:this.handleBlur,onKeyup:this.handleKeyup,onKeydown:this.handleKeydown},[K("div",{class:w(`${t}-switch__rail`),"aria-hidden":"true",style:ve(d)},[B(()=>O(f,s=>O(S,l=>s||l?(p(),D("div",{key:4,"aria-hidden":!0,class:w(`${t}-switch__children-placeholder`)},[K("div",{class:w(`${t}-switch__rail-placeholder`)},[K("div",{class:w(`${t}-switch__button-placeholder`)},null,2),B(()=>s)],2),K("div",{class:w(`${t}-switch__rail-placeholder`)},[K("div",{class:w(`${t}-switch__button-placeholder`)},null,2),B(()=>l)],2)],2)):null))),K("div",{class:w(`${t}-switch__button`)},[B(()=>O(b,s=>O(C,l=>O(c,z=>(p(),j(Me,null,{default:()=>this.loading?(p(),j(Pe,We({key:"loading",clsPrefix:t,strokeWidth:20},this.spinProps),null,16,["clsPrefix"])):this.checked&&(l||s)?(p(),D("div",{class:w(`${t}-switch__button-icon`),key:l?"checked-icon":"icon"},[B(()=>l||s)],2)):!this.checked&&(z||s)?(p(),D("div",{class:w(`${t}-switch__button-icon`),key:z?"unchecked-icon":"icon"},[B(()=>z||s)],2)):null},1024)))))),B(()=>O(f,s=>s&&(p(),D("div",{key:"checked",class:w(`${t}-switch__checked`)},[B(()=>s)],2)))),B(()=>O(S,s=>s&&(p(),D("div",{key:"unchecked",class:w(`${t}-switch__unchecked`)},[B(()=>s)],2))))],2)],6)],46,oa)}});const ba=me({__name:"Users",setup(t){const h=Je(),y=v(!1),d=v(!1),k=v([]),x=v(0),f=de({username:"",page:1,page_size:10}),S=v([]),b=v(!1),C=v(!1),c=v(!1),_=v(0),s=v(""),l=de({username:"",password:"",nickname:"",email:"",role_ids:[],statusOn:1}),z=de({page:1,pageSize:10,pageCount:1,showSizePicker:!0,onChange:i=>{f.page=i,m()},onUpdatePageSize:i=>{f.page_size=i,m()}});async function m(){y.value=!0;try{const{data:i}=await je(f);k.value=i.data.list,x.value=i.data.page.total,z.page=f.page,z.pageSize=f.page_size,z.pageCount=Math.ceil(x.value/f.page_size)}finally{y.value=!1}}async function Y(){const{data:i}=await He({page:1,page_size:100});S.value=i.data.list.map(a=>({label:a.name,value:a.id}))}function q(){c.value=!1,Object.assign(l,{username:"",password:"",nickname:"",email:"",role_ids:[],statusOn:1}),b.value=!0}function G(i){c.value=!0,_.value=i.id,Object.assign(l,{username:i.username,password:"",nickname:i.nickname,email:i.email,role_ids:i.role_ids??[],statusOn:i.status}),b.value=!0}async function Q(){var i,a;d.value=!0;try{c.value?(await Ae(_.value,{nickname:l.nickname,email:l.email,status:l.statusOn}),await Ee(_.value,l.role_ids)):await Xe({username:l.username,password:l.password,nickname:l.nickname,email:l.email,role_ids:l.role_ids}),h.success("保存成功"),b.value=!1,await m()}catch(g){h.error(((a=(i=g==null?void 0:g.response)==null?void 0:i.data)==null?void 0:a.msg)||"保存失败")}finally{d.value=!1}}async function Z(){if(s.value.length<6){h.warning("密码至少 6 位");return}await Ye(_.value,s.value),h.success("密码已重置"),C.value=!1}const J=[{title:"ID",key:"id",width:60},{title:"用户名",key:"username"},{title:"昵称",key:"nickname"},{title:"邮箱",key:"email"},{title:"状态",key:"status",width:80,render:i=>A(Ze,{type:i.status===1?"success":"error",size:"small"},{default:()=>i.status===1?"启用":"禁用"})},{title:"创建时间",key:"created_at",width:170},{title:"操作",key:"actions",width:220,render(i){return A(be,{},{default:()=>[A(V,{size:"small",onClick:()=>G(i)},{default:()=>"编辑"}),A(V,{size:"small",type:"warning",onClick:()=>{_.value=i.id,s.value="",C.value=!0}},{default:()=>"重置密码"}),A(V,{size:"small",type:"error",onClick:async()=>{await qe(i.id),h.success("已删除"),m()}},{default:()=>"删除"})]})}}];return Te(()=>{m(),Y()}),(i,a)=>{const g=De("permission");return p(),D(Ie,null,[r(n(Ne),{title:"用户管理"},{"header-extra":u(()=>[r(n(be),null,{default:u(()=>[r(n(N),{value:f.username,"onUpdate:value":a[0]||(a[0]=e=>f.username=e),placeholder:"用户名",clearable:"",style:{width:"160px"},onKeyup:Ke(m,["enter"])},null,8,["value"]),r(n(V),{type:"primary",onClick:m},{default:u(()=>[...a[12]||(a[12]=[L("查询",-1)])]),_:1}),Le((p(),j(n(V),{type:"primary",ghost:"",onClick:q},{default:u(()=>[...a[13]||(a[13]=[L("新增用户",-1)])]),_:1})),[[g,["menu:user"]]])]),_:1})]),default:u(()=>[r(n(Ge),{columns:J,data:k.value,loading:y.value,pagination:z,remote:""},null,8,["data","loading","pagination"])]),_:1}),r(n(pe),{show:b.value,"onUpdate:show":a[8]||(a[8]=e=>b.value=e),preset:"dialog",title:c.value?"编辑用户":"新增用户",style:{width:"480px"}},{action:u(()=>[r(n(V),{onClick:a[7]||(a[7]=e=>b.value=!1)},{default:u(()=>[...a[14]||(a[14]=[L("取消",-1)])]),_:1}),r(n(V),{type:"primary",loading:d.value,onClick:Q},{default:u(()=>[...a[15]||(a[15]=[L("确定",-1)])]),_:1},8,["loading"])]),default:u(()=>[r(n(ea),{model:l,"label-placement":"left","label-width":"80"},{default:u(()=>[c.value?ue("",!0):(p(),j(n(I),{key:0,label:"用户名"},{default:u(()=>[r(n(N),{value:l.username,"onUpdate:value":a[1]||(a[1]=e=>l.username=e)},null,8,["value"])]),_:1})),c.value?ue("",!0):(p(),j(n(I),{key:1,label:"密码"},{default:u(()=>[r(n(N),{value:l.password,"onUpdate:value":a[2]||(a[2]=e=>l.password=e),type:"password"},null,8,["value"])]),_:1})),r(n(I),{label:"昵称"},{default:u(()=>[r(n(N),{value:l.nickname,"onUpdate:value":a[3]||(a[3]=e=>l.nickname=e)},null,8,["value"])]),_:1}),r(n(I),{label:"邮箱"},{default:u(()=>[r(n(N),{value:l.email,"onUpdate:value":a[4]||(a[4]=e=>l.email=e)},null,8,["value"])]),_:1}),r(n(I),{label:"角色"},{default:u(()=>[r(n(Qe),{value:l.role_ids,"onUpdate:value":a[5]||(a[5]=e=>l.role_ids=e),multiple:"",options:S.value},null,8,["value","options"])]),_:1}),c.value?(p(),j(n(I),{key:2,label:"状态"},{default:u(()=>[r(n(ra),{value:l.statusOn,"onUpdate:value":a[6]||(a[6]=e=>l.statusOn=e),"checked-value":1,"unchecked-value":0},null,8,["value"])]),_:1})):ue("",!0)]),_:1},8,["model"])]),_:1},8,["show","title"]),r(n(pe),{show:C.value,"onUpdate:show":a[11]||(a[11]=e=>C.value=e),preset:"dialog",title:"重置密码",style:{width:"420px"}},{action:u(()=>[r(n(V),{onClick:a[10]||(a[10]=e=>C.value=!1)},{default:u(()=>[...a[16]||(a[16]=[L("取消",-1)])]),_:1}),r(n(V),{type:"primary",onClick:Z},{default:u(()=>[...a[17]||(a[17]=[L("确定",-1)])]),_:1})]),default:u(()=>[r(n(N),{value:s.value,"onUpdate:value":a[9]||(a[9]=e=>s.value=e),type:"password",placeholder:"新密码（至少 6 位）"},null,8,["value"])]),_:1},8,["show"])],64)}}});export{ba as default};
