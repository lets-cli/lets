(window.webpackJsonp=window.webpackJsonp||[]).push([[15],{105:function(e,t,n){"use strict";n.d(t,"a",(function(){return b})),n.d(t,"b",(function(){return m}));var r=n(0),a=n.n(r);function l(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function o(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){l(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},l=Object.keys(e);for(r=0;r<l.length;r++)n=l[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var l=Object.getOwnPropertySymbols(e);for(r=0;r<l.length;r++)n=l[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var c=a.a.createContext({}),p=function(e){var t=a.a.useContext(c),n=t;return e&&(n="function"==typeof e?e(t):o(o({},t),e)),n},b=function(e){var t=p(e.components);return a.a.createElement(c.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},d=a.a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,l=e.originalType,i=e.parentName,c=s(e,["components","mdxType","originalType","parentName"]),b=p(n),d=r,m=b["".concat(i,".").concat(d)]||b[d]||u[d]||l;return n?a.a.createElement(m,o(o({ref:t},c),{},{components:n})):a.a.createElement(m,o({ref:t},c))}));function m(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var l=n.length,i=new Array(l);i[0]=d;var o={};for(var s in t)hasOwnProperty.call(t,s)&&(o[s]=t[s]);o.originalType=e,o.mdxType="string"==typeof e?e:r,i[1]=o;for(var c=2;c<l;c++)i[c]=n[c];return a.a.createElement.apply(null,i)}return a.a.createElement.apply(null,n)}d.displayName="MDXCreateElement"},86:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return i})),n.d(t,"metadata",(function(){return o})),n.d(t,"toc",(function(){return s})),n.d(t,"default",(function(){return p}));var r=n(3),a=n(7),l=(n(0),n(105)),i={id:"development",title:"Development"},o={unversionedId:"development",id:"development",isDocsHomePage:!1,title:"Development",description:"Build",source:"@site/docs/development.md",slug:"/development",permalink:"/docs/development",editUrl:"https://github.com/lets-cli/lets/edit/master/docs/docs/development.md",version:"current",sidebar:"someSidebar",previous:{title:"IDE/Text editors support",permalink:"/docs/ide_support"},next:{title:"Contribute",permalink:"/docs/contribute"}},s=[{value:"Build",id:"build",children:[]},{value:"Test",id:"test",children:[]},{value:"Release",id:"release",children:[]},{value:"Versioning",id:"versioning",children:[]}],c={toc:s};function p(e){var t=e.components,n=Object(a.a)(e,["components"]);return Object(l.b)("wrapper",Object(r.a)({},c,n,{components:t,mdxType:"MDXLayout"}),Object(l.b)("h2",{id:"build"},"Build"),Object(l.b)("p",null,"To build a binary:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre",className:"language-bash"},"go build -o lets *.go\n")),Object(l.b)("p",null,"To install in system"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre",className:"language-bash"},"go build -o lets *.go && sudo mv ./lets /usr/local/bin/lets\n")),Object(l.b)("p",null,"Or if you already have ",Object(l.b)("inlineCode",{parentName:"p"},"lets")," installed in your system:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre",className:"language-bash"},"lets build-and-install\n")),Object(l.b)("p",null,"After install - check version of lets - ",Object(l.b)("inlineCode",{parentName:"p"},"lets --version")," - it should be development"),Object(l.b)("p",null,"It will install ",Object(l.b)("inlineCode",{parentName:"p"},"lets")," to /usr/local/bin/lets and set version to development with current tag and timestamp"),Object(l.b)("h2",{id:"test"},"Test"),Object(l.b)("p",null,"To run all tests:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre",className:"language-shell",metastring:"script",script:!0},"lets test\n")),Object(l.b)("p",null,"To run unit tests:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre",className:"language-shell",metastring:"script",script:!0},"lets test-unit\n")),Object(l.b)("p",null,"To get coverage:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre",className:"language-shell",metastring:"script",script:!0},"lets coverage\n")),Object(l.b)("p",null,"To test ",Object(l.b)("inlineCode",{parentName:"p"},"lets")," output we using ",Object(l.b)("inlineCode",{parentName:"p"},"bats")," - bash automated testing:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre",className:"language-shell",metastring:"script",script:!0},"lets test-bats\n\n# or run one test\n\nlets test-bats global_env.bats\n")),Object(l.b)("h2",{id:"release"},"Release"),Object(l.b)("p",null,"To release a new version:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre",className:"language-bash"},'lets release 0.0.1 -m "implement some new feature"\n')),Object(l.b)("p",null,"This will create an annotated tag with 0.0.1 and run ",Object(l.b)("inlineCode",{parentName:"p"},"git push --tags")),Object(l.b)("h2",{id:"versioning"},"Versioning"),Object(l.b)("p",null,Object(l.b)("inlineCode",{parentName:"p"},"lets")," releases must be backward compatible. That means every new ",Object(l.b)("inlineCode",{parentName:"p"},"lets")," release must work with old configs."),Object(l.b)("p",null,"For situations like e.g. new functionality, there is a ",Object(l.b)("inlineCode",{parentName:"p"},"version")," in ",Object(l.b)("inlineCode",{parentName:"p"},"lets.yaml")," which specifies ",Object(l.b)("strong",{parentName:"p"},"minimum required")," ",Object(l.b)("inlineCode",{parentName:"p"},"lets")," version."),Object(l.b)("p",null,"If ",Object(l.b)("inlineCode",{parentName:"p"},"lets")," version installed on the user machine is less than the one specified in config it will show and error and ask the user to upgrade ",Object(l.b)("inlineCode",{parentName:"p"},"lets")," version."))}p.isMDXComponent=!0}}]);