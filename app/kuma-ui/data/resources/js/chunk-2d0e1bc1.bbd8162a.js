(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-2d0e1bc1"],{"7c65":function(t,n,e){"use strict";e.r(n);var o=function(){var t=this,n=t.$createElement,e=t._self._c||n;return e("div",{staticClass:"dataplanes-detail"},[e("YamlView",{attrs:{title:"Entity Overview","has-error":t.hasError,"is-loading":t.isLoading,"is-empty":t.isEmpty,content:t.content}})],1)},i=[],a=e("ff9d"),r={name:"TrafficRouteDetail",metaInfo:{title:"Traffic Route Details"},components:{YamlView:a["a"]},data:function(){return{content:null,hasError:!1,isLoading:!0,isEmpty:!1}},watch:{$route:function(t,n){this.bootstrap()}},beforeMount:function(){this.bootstrap()},methods:{bootstrap:function(){var t=this,n=this.$route.params.mesh,e=this.$route.params.trafficroute;return this.$api.getTrafficRoute(n,e).then((function(n){n?t.content=n:t.$router.push("/404")})).catch((function(n){t.hasError=!0,console.error(n)})).finally((function(){setTimeout((function(){t.isLoading=!1}),"500")}))}}},s=r,c=e("2877"),u=Object(c["a"])(s,o,i,!1,null,null,null);n["default"]=u.exports}}]);