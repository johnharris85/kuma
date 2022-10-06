import{W,X as F,J as R,I as G,C as H,O as X,ab as Z,K as B,_ as U,r as D,o as e,c as g,w as l,a as d,e as n,h as v,f as a,n as A,t as c,F as _,b as E,d as Y,i as x,j as w,m as p,cF as J,cG as Q,cL as ee,cA as M,cH as ae,a6 as te,q as j,cB as se,cD as ne,s as oe,cM as le,p as ie,cN as re,a5 as ce,x as de,y as ue,u as pe,k as V}from"./index.8c6a97c0.js";import{S as z,E as $}from"./EnvoyData.5dfd7178.js";import{E as me}from"./EntityURLControl.a9e8b804.js";import{L as q}from"./LabelList.f7f25a11.js";import{a as he,S as _e}from"./SubscriptionHeader.3ae591ed.js";import{T as ye}from"./TabsWidget.10ca7693.js";import{W as fe}from"./WarningsWidget.a388dc74.js";import{Y as ve}from"./YamlView.b86fd036.js";import{_ as ge}from"./EmptyBlock.vue_vue_type_script_setup_true_lang.80be7515.js";import{E as Pe}from"./ErrorBlock.d4c6dfc2.js";import{_ as be}from"./LoadingBlock.vue_vue_type_script_setup_true_lang.1c68a0d8.js";import"./CodeBlock.31e0047b.js";import"./_commonjsHelpers.712cc82f.js";import"./index.58caa11d.js";const ke={inbound:"Policies applied on incoming connection on address",outbound:"Policies applied on outgoing connection to the address",service:"Policies applied on outgoing connections to service",dataplane:"Policies applied on all incoming and outgoing connections to the selected data plane proxy"},De={name:"DataplanePolicies",components:{StatusInfo:z,AccordionList:W,AccordionItem:F,KCard:R,KPop:G,KIcon:H,KBadge:X},props:{mesh:{type:String,required:!0},dppName:{type:String,required:!0}},data(){return{items:[],hasItems:!1,isLoading:!0,error:null,searchInput:"",POLICY_TYPE_SUBTITLE:ke}},computed:{...Z({policiesByType:s=>s.policiesByType})},watch:{dppName(){this.fetchPolicies()}},mounted(){this.fetchPolicies()},methods:{async fetchPolicies(){this.error=null,this.isLoading=!0;try{const{items:s,total:o,kind:h}=await B.getDataplanePolicies({mesh:this.mesh,dppName:this.dppName});(h===void 0||h==="SidecarDataplane")&&(this.processItems(s),this.items=s,this.hasItems=o>0)}catch(s){this.error=s}finally{this.isLoading=!1}},processItems(s){for(const o of s){o.policyTypes={};for(const h in o.matchedPolicies){const y=this.policiesByType[h],u={pluralDisplayName:y.pluralDisplayName,policies:o.matchedPolicies[h]};for(const f of u.policies)f.route={name:y.path,query:{ns:f.name},params:{mesh:f.mesh}};o.policyTypes[h]=u}}}}},we={class:"flex items-center justify-between"},Ie={key:0,class:"text-lg"},Te={key:1,class:"text-lg"},Se={class:"subtitle"},Ce={key:0},Le={key:1},Oe={class:"flex flex-wrap justify-end"},xe={class:"policy-wrapper"},Ae={class:"policy-type"};function Ee(s,o,h,y,u,f){const P=D("KIcon"),S=D("KPop"),I=D("KBadge"),b=D("router-link"),k=D("AccordionItem"),C=D("AccordionList"),T=D("KCard"),N=D("StatusInfo");return e(),g(N,{"has-error":u.error!==null,"is-loading":u.isLoading,error:u.error,"is-empty":!u.hasItems},{default:l(()=>[d(T,{"border-variant":"noBorder"},{body:l(()=>[d(C,{"initially-open":[],"multiple-open":""},{default:l(()=>[(e(!0),n(_,null,v(u.items,(t,m)=>(e(),g(k,{key:m},{"accordion-header":l(()=>[a("div",we,[a("div",null,[t.type==="dataplane"?(e(),n("p",Ie," Dataplane ")):A("",!0),t.type!=="dataplane"?(e(),n("p",Te,c(t.service),1)):A("",!0),a("p",Se,[t.type==="inbound"||t.type==="outbound"?(e(),n("span",Ce,c(t.type)+" "+c(t.name),1)):t.type==="service"||t.type==="dataplane"?(e(),n("span",Le,c(t.type),1)):A("",!0),d(S,{width:"300",placement:"right",trigger:"hover"},{content:l(()=>[a("div",null,c(u.POLICY_TYPE_SUBTITLE[t.type]),1)]),default:l(()=>[d(P,{icon:"help",size:"16",class:"ml-1"})]),_:2},1024)])]),a("div",Oe,[(e(!0),n(_,null,v(t.matchedPolicies,(r,i)=>(e(),g(I,{key:`${m}-${i}`,class:"mr-2 mb-2"},{default:l(()=>[E(c(i),1)]),_:2},1024))),128))])])]),"accordion-content":l(()=>[a("div",xe,[(e(!0),n(_,null,v(t.policyTypes,(r,i)=>(e(),n("div",{key:`${m}-${i}`,class:"policy-item"},[a("h4",Ae,c(r.pluralDisplayName),1),a("ul",null,[(e(!0),n(_,null,v(r.policies,(L,O)=>(e(),n("li",{key:`${m}-${i}-${O}`,class:"my-1","data-testid":"policy-name"},[d(b,{to:L.route},{default:l(()=>[E(c(L.name),1)]),_:2},1032,["to"])]))),128))])]))),128))])]),_:2},1024))),128))]),_:1})]),_:1})]),_:1},8,["has-error","is-loading","error","is-empty"])}const Ne=U(De,[["render",Ee],["__scopeId","data-v-7223f9a8"]]),K=s=>(de("data-v-5f6db0c0"),s=s(),ue(),s),$e={key:0},Be={class:"entity-status__label"},Ke=K(()=>a("span",{class:"entity-status__dot"},null,-1)),Me={key:1},Ve=K(()=>a("h4",null,"Tags",-1)),qe=K(()=>a("h4",null,"Versions",-1)),We={class:"config-wrapper"},Fe={key:0},Re=["href"],Ue=Y({__name:"DataPlaneDetails",props:{dataPlane:{type:Object,required:!0},dataPlaneOverview:{type:Object,required:!0}},setup(s){const o=s,h=j(),y=[{hash:"#overview",title:"Overview"},{hash:"#insights",title:"DPP Insights"},{hash:"#dpp-policies",title:"Policies"},{hash:"#xds-configuration",title:"XDS Configuration"},{hash:"#envoy-stats",title:"Stats"},{hash:"#envoy-clusters",title:"Clusters"},{hash:"#mtls",title:"Certificate Insights"},{hash:"#warnings",title:"Warnings"}],u=x([]),f=w(()=>{const{type:t,name:m,mesh:r}=o.dataPlane,i=se(o.dataPlane,o.dataPlaneOverview.dataplaneInsight);return{type:t,name:m,mesh:r,status:i}}),P=w(()=>M(o.dataPlane)),S=w(()=>ne(o.dataPlaneOverview.dataplaneInsight)),I=w(()=>oe(o.dataPlane)),b=w(()=>le(o.dataPlaneOverview)),k=w(()=>{const t=Array.from(o.dataPlaneOverview.dataplaneInsight.subscriptions);return t.reverse(),t}),C=w(()=>{const t=h.getters.getKumaDocsVersion;return t!==null?t:"latest"}),T=w(()=>u.value.length===0?y.filter(t=>t.hash!=="#warnings"):y);function N(){const t=o.dataPlaneOverview.dataplaneInsight.subscriptions;if(t.length===0||!("version"in t[0]))return;const m=t[0].version;if(m.kumaDp&&m.envoy){const i=J(m);i.kind!==Q&&i.kind!==ee&&u.value.push(i)}h.getters["config/getMulticlusterStatus"]&&M(o.dataPlane).find(O=>O.label===ae)&&typeof m.kumaDp.kumaCpCompatible=="boolean"&&!m.kumaDp.kumaCpCompatible&&u.value.push({kind:te,payload:{kumaDp:m.kumaDp.version}})}return N(),(t,m)=>(e(),g(ye,{tabs:p(T),"initial-tab-override":"overview"},{tabHeader:l(()=>[a("div",null,[a("h3",null," DPP: "+c(s.dataPlane.name),1)]),a("div",null,[d(me,{name:s.dataPlane.name,mesh:s.dataPlane.mesh},null,8,["name","mesh"])])]),overview:l(()=>[d(q,null,{default:l(()=>[a("div",null,[a("ul",null,[(e(!0),n(_,null,v(p(f),(r,i)=>(e(),n("li",{key:i},[a("h4",null,c(i),1),i==="status"&&typeof r!="string"?(e(),n("div",$e,[a("div",{class:ie(["entity-status",{"is-offline":r.status.toString().toLowerCase()==="offline","is-degraded":r.status.toString().toLowerCase()==="partially degraded"}])},[a("span",Be,c(r.status),1)],2),(e(!0),n(_,null,v(r.reason,(L,O)=>(e(),n("div",{key:O,class:"reason"},[Ke,E(" "+c(L),1)]))),128))])):(e(),n("div",Me,c(r),1))]))),128))])]),a("div",null,[p(P).length>0?(e(),n(_,{key:0},[Ve,a("p",null,[(e(!0),n(_,null,v(p(P),(r,i)=>(e(),n("span",{key:i,class:"tag-cols"},[a("span",null,c(r.label)+": ",1),a("span",null,c(r.value),1)]))),128))])],64)):A("",!0),p(S)?(e(),n(_,{key:1},[qe,a("p",null,[(e(!0),n(_,null,v(p(S),(r,i)=>(e(),n("span",{key:i,class:"tag-cols"},[a("span",null,c(i)+": ",1),a("span",null,c(r),1)]))),128))])],64)):A("",!0)])]),_:1}),a("div",We,[d(ve,{id:"code-block-data-plane",content:p(I),"is-searchable":""},null,8,["content"])])]),insights:l(()=>[d(z,{"is-empty":p(k).length===0},{default:l(()=>[d(p(R),{"border-variant":"noBorder"},{body:l(()=>[d(W,{"initially-open":0},{default:l(()=>[(e(!0),n(_,null,v(p(k),(r,i)=>(e(),g(F,{key:i},{"accordion-header":l(()=>[d(he,{details:r},null,8,["details"])]),"accordion-content":l(()=>[d(_e,{details:r,"is-discovery-subscription":""},null,8,["details"])]),_:2},1024))),128))]),_:1})]),_:1})]),_:1},8,["is-empty"])]),"dpp-policies":l(()=>[d(Ne,{mesh:s.dataPlane.mesh,"dpp-name":s.dataPlane.name},null,8,["mesh","dpp-name"])]),"xds-configuration":l(()=>[d($,{"data-path":"xds",mesh:s.dataPlane.mesh,"dpp-name":s.dataPlane.name,"query-key":"envoy-data-data-plane"},null,8,["mesh","dpp-name"])]),"envoy-stats":l(()=>[d($,{"data-path":"stats",mesh:s.dataPlane.mesh,"dpp-name":s.dataPlane.name,"query-key":"envoy-data-data-plane"},null,8,["mesh","dpp-name"])]),"envoy-clusters":l(()=>[d($,{"data-path":"clusters",mesh:s.dataPlane.mesh,"dpp-name":s.dataPlane.name,"query-key":"envoy-data-data-plane"},null,8,["mesh","dpp-name"])]),mtls:l(()=>[d(q,null,{default:l(()=>[p(b)!==null?(e(),n("ul",Fe,[(e(!0),n(_,null,v(p(b),(r,i)=>(e(),n("li",{key:i},[a("h4",null,c(r.label),1),a("p",null,c(r.value),1)]))),128))])):(e(),g(p(re),{key:1,appearance:"danger"},{alertMessage:l(()=>[E(" This data plane proxy does not yet have mTLS configured \u2014 "),a("a",{href:`https://kuma.io/docs/${p(C)}/documentation/security/#certificates`,class:"external-link",target:"_blank"}," Learn About Certificates in "+c(p(ce)),9,Re)]),_:1}))]),_:1})]),warnings:l(()=>[d(fe,{warnings:u.value},null,8,["warnings"])]),_:1},8,["tabs"]))}});const Ye=U(Ue,[["__scopeId","data-v-5f6db0c0"]]),je={class:"component-frame"},ia=Y({__name:"DataPlaneDetailView",setup(s){const o=pe(),h=j(),y=x(null),u=x(null),f=x(!0),P=x(null),S={};async function I(){P.value=null,f.value=!0;const b=o.params.mesh,k=o.params.dataPlane,C=S;try{y.value=await B.getDataplaneFromMesh({mesh:b,name:k},C),u.value=await B.getDataplaneOverviewFromMesh({mesh:b,name:k},C)}catch(T){y.value=null,T instanceof Error?P.value=T:console.error(T)}finally{f.value=!1}}return V(()=>o.params.mesh,function(){o.name==="data-plane-detail-view"&&I()}),V(()=>o.params.dataPlane,function(){o.name==="data-plane-detail-view"&&I()}),I(),h.dispatch("updatePageTitle",o.params.dataPlane),(b,k)=>(e(),n("div",je,[f.value?(e(),g(be,{key:0})):P.value!==null?(e(),g(Pe,{key:1,error:P.value},null,8,["error"])):y.value===null||u.value===null?(e(),g(ge,{key:2})):(e(),g(Ye,{key:3,"data-plane":y.value,"data-plane-overview":u.value},null,8,["data-plane","data-plane-overview"]))]))}});export{ia as default};
