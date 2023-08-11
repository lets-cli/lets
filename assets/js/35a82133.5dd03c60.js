"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[83],{3905:function(e,n,t){t.d(n,{Zo:function(){return c},kt:function(){return d}});var r=t(7294);function a(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function i(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);n&&(r=r.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,r)}return t}function l(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?i(Object(t),!0).forEach((function(n){a(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):i(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function o(e,n){if(null==e)return{};var t,r,a=function(e,n){if(null==e)return{};var t,r,a={},i=Object.keys(e);for(r=0;r<i.length;r++)t=i[r],n.indexOf(t)>=0||(a[t]=e[t]);return a}(e,n);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)t=i[r],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(a[t]=e[t])}return a}var p=r.createContext({}),m=function(e){var n=r.useContext(p),t=n;return e&&(t="function"==typeof e?e(n):l(l({},n),e)),t},c=function(e){var n=m(e.components);return r.createElement(p.Provider,{value:n},e.children)},s={inlineCode:"code",wrapper:function(e){var n=e.children;return r.createElement(r.Fragment,{},n)}},u=r.forwardRef((function(e,n){var t=e.components,a=e.mdxType,i=e.originalType,p=e.parentName,c=o(e,["components","mdxType","originalType","parentName"]),u=m(t),d=a,f=u["".concat(p,".").concat(d)]||u[d]||s[d]||i;return t?r.createElement(f,l(l({ref:n},c),{},{components:t})):r.createElement(f,l({ref:n},c))}));function d(e,n){var t=arguments,a=n&&n.mdxType;if("string"==typeof e||a){var i=t.length,l=new Array(i);l[0]=u;var o={};for(var p in n)hasOwnProperty.call(n,p)&&(o[p]=n[p]);o.originalType=e,o.mdxType="string"==typeof e?e:a,l[1]=o;for(var m=2;m<i;m++)l[m]=t[m];return r.createElement.apply(null,l)}return r.createElement.apply(null,t)}u.displayName="MDXCreateElement"},5791:function(e,n,t){t.r(n),t.d(n,{assets:function(){return c},contentTitle:function(){return p},default:function(){return d},frontMatter:function(){return o},metadata:function(){return m},toc:function(){return s}});var r=t(7462),a=t(3366),i=(t(7294),t(3905)),l=["components"],o={id:"env",title:"Environment"},p=void 0,m={unversionedId:"env",id:"env",title:"Environment",description:"Default environment variables",source:"@site/docs/env.md",sourceDirName:".",slug:"/env",permalink:"/docs/env",draft:!1,editUrl:"https://github.com/lets-cli/lets/edit/master/docs/docs/env.md",tags:[],version:"current",frontMatter:{id:"env",title:"Environment"},sidebar:"someSidebar",previous:{title:"CLI options",permalink:"/docs/cli"},next:{title:"Examples",permalink:"/docs/examples"}},c={},s=[{value:"Default environment variables",id:"default-environment-variables",level:3},{value:"Environment variables available at command runtime",id:"environment-variables-available-at-command-runtime",level:3},{value:"Override command env with -E flag",id:"override-command-env-with--e-flag",level:3}],u={toc:s};function d(e){var n=e.components,t=(0,a.Z)(e,l);return(0,i.kt)("wrapper",(0,r.Z)({},u,t,{components:n,mdxType:"MDXLayout"}),(0,i.kt)("h3",{id:"default-environment-variables"},"Default environment variables"),(0,i.kt)("p",null,(0,i.kt)("inlineCode",{parentName:"p"},"lets")," has builtin environ variables which user can override before lets execution. E.g ",(0,i.kt)("inlineCode",{parentName:"p"},"LETS_DEBUG=1 lets test")),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"LETS_DEBUG")," - enable debug messages"),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"LETS_CONFIG")," - changes default ",(0,i.kt)("inlineCode",{parentName:"li"},"lets.yaml")," file path (e.g. LETS_CONFIG=lets.my.yaml)"),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"LETS_CONFIG_DIR")," - changes path to dir where ",(0,i.kt)("inlineCode",{parentName:"li"},"lets.yaml")," file placed"),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"NO_COLOR")," - disables colored output. See ",(0,i.kt)("a",{parentName:"li",href:"https://no-color.org/"},"https://no-color.org/"))),(0,i.kt)("h3",{id:"environment-variables-available-at-command-runtime"},"Environment variables available at command runtime"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"LETS_COMMAND_NAME")," - string name of launched command"),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"LETS_COMMAND_ARGS")," - positional arguments for launched command, e.g. for ",(0,i.kt)("inlineCode",{parentName:"li"},"lets run --debug --config=test.ini")," it will contain ",(0,i.kt)("inlineCode",{parentName:"li"},"--debug --config=test.ini")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"LETS_COMMAND_WORK_DIR")," - absolute path to ",(0,i.kt)("inlineCode",{parentName:"li"},"work_dir")," specified in command."),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"LETS_CONFIG")," - absolute path to lets config file."),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"LETS_CONFIG_DIR")," - absolute path to lets config file firectory."),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"LETS_SHELL")," - shell from config or command."),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"LETSOPT_<>")," - options parsed from command ",(0,i.kt)("inlineCode",{parentName:"li"},"options")," (docopt string). E.g ",(0,i.kt)("inlineCode",{parentName:"li"},"lets run --env=prod --reload")," will be ",(0,i.kt)("inlineCode",{parentName:"li"},"LETSOPT_ENV=prod")," and ",(0,i.kt)("inlineCode",{parentName:"li"},"LETSOPT_RELOAD=true")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("inlineCode",{parentName:"li"},"LETSCLI_<>")," - options which values is a options usage. E.g ",(0,i.kt)("inlineCode",{parentName:"li"},"lets run --env=prod --reload")," will be ",(0,i.kt)("inlineCode",{parentName:"li"},"LETSCLI_ENV=--env=prod")," and ",(0,i.kt)("inlineCode",{parentName:"li"},"LETSCLI_RELOAD=--reload"))),(0,i.kt)("h3",{id:"override-command-env-with--e-flag"},"Override command env with -E flag"),(0,i.kt)("p",null,"You can override environment for command with ",(0,i.kt)("inlineCode",{parentName:"p"},"-E")," flag:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-yaml"},"shell: bash\n\ncommands:\n  say:\n    env:\n      NAME: Rick\n    cmd: echo Hello ${NAME}\n")),(0,i.kt)("p",null,(0,i.kt)("strong",{parentName:"p"},(0,i.kt)("inlineCode",{parentName:"strong"},"lets say"))," - prints ",(0,i.kt)("inlineCode",{parentName:"p"},"Hello Rick")),(0,i.kt)("p",null,(0,i.kt)("strong",{parentName:"p"},(0,i.kt)("inlineCode",{parentName:"strong"},"lets -E NAME=Morty say"))," - prints ",(0,i.kt)("inlineCode",{parentName:"p"},"Hello Morty")),(0,i.kt)("p",null,"Alternatively:"),(0,i.kt)("p",null,(0,i.kt)("strong",{parentName:"p"},(0,i.kt)("inlineCode",{parentName:"strong"},"lets --env NAME=Morty say"))," - prints ",(0,i.kt)("inlineCode",{parentName:"p"},"Hello Morty")))}d.isMDXComponent=!0}}]);