(window.webpackJsonp=window.webpackJsonp||[]).push([[23],{107:function(e,n,t){"use strict";t.d(n,"a",(function(){return d})),t.d(n,"b",(function(){return O}));var a=t(0),i=t.n(a);function l(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function r(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);n&&(a=a.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,a)}return t}function c(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?r(Object(t),!0).forEach((function(n){l(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):r(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function b(e,n){if(null==e)return{};var t,a,i=function(e,n){if(null==e)return{};var t,a,i={},l=Object.keys(e);for(a=0;a<l.length;a++)t=l[a],n.indexOf(t)>=0||(i[t]=e[t]);return i}(e,n);if(Object.getOwnPropertySymbols){var l=Object.getOwnPropertySymbols(e);for(a=0;a<l.length;a++)t=l[a],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(i[t]=e[t])}return i}var o=i.a.createContext({}),p=function(e){var n=i.a.useContext(o),t=n;return e&&(t="function"==typeof e?e(n):c(c({},n),e)),t},d=function(e){var n=p(e.components);return i.a.createElement(o.Provider,{value:n},e.children)},m={inlineCode:"code",wrapper:function(e){var n=e.children;return i.a.createElement(i.a.Fragment,{},n)}},s=i.a.forwardRef((function(e,n){var t=e.components,a=e.mdxType,l=e.originalType,r=e.parentName,o=b(e,["components","mdxType","originalType","parentName"]),d=p(t),s=a,O=d["".concat(r,".").concat(s)]||d[s]||m[s]||l;return t?i.a.createElement(O,c(c({ref:n},o),{},{components:t})):i.a.createElement(O,c({ref:n},o))}));function O(e,n){var t=arguments,a=n&&n.mdxType;if("string"==typeof e||a){var l=t.length,r=new Array(l);r[0]=s;var c={};for(var b in n)hasOwnProperty.call(n,b)&&(c[b]=n[b]);c.originalType=e,c.mdxType="string"==typeof e?e:a,r[1]=c;for(var o=2;o<l;o++)r[o]=t[o];return i.a.createElement.apply(null,r)}return i.a.createElement.apply(null,t)}s.displayName="MDXCreateElement"},94:function(e,n,t){"use strict";t.r(n),t.d(n,"frontMatter",(function(){return r})),t.d(n,"metadata",(function(){return c})),t.d(n,"toc",(function(){return b})),t.d(n,"default",(function(){return p}));var a=t(3),i=t(7),l=(t(0),t(107)),r={id:"changelog",title:"Changelog"},c={unversionedId:"changelog",id:"changelog",isDocsHomePage:!1,title:"Changelog",description:"[Unreleased]",source:"@site/docs/changelog.md",slug:"/changelog",permalink:"/docs/changelog",editUrl:"https://github.com/lets-cli/lets/edit/master/docs/docs/changelog.md",version:"current",sidebar:"someSidebar",previous:{title:"Best practices",permalink:"/docs/best_practices"},next:{title:"IDE/Text editors support",permalink:"/docs/ide_support"}},b=[{value:"Unreleased",id:"unreleased",children:[]},{value:"0.0.49",id:"0049",children:[]},{value:"0.0.48",id:"0048",children:[]},{value:"0.0.47",id:"0047",children:[]},{value:"0.0.45",id:"0045",children:[]},{value:"0.0.44",id:"0044",children:[]},{value:"0.0.43",id:"0043",children:[]},{value:"0.0.42",id:"0042",children:[]},{value:"0.0.41",id:"0041",children:[]},{value:"0.0.40",id:"0040",children:[]},{value:"0.0.33",id:"0033",children:[]},{value:"0.0.32",id:"0032",children:[]},{value:"0.0.30",id:"0030",children:[]},{value:"0.0.29",id:"0029",children:[]},{value:"0.0.28",id:"0028",children:[]},{value:"0.0.27",id:"0027",children:[]}],o={toc:b};function p(e){var n=e.components,t=Object(i.a)(e,["components"]);return Object(l.b)("wrapper",Object(a.a)({},o,t,{components:n,mdxType:"MDXLayout"}),Object(l.b)("h2",{id:"unreleased"},"[Unreleased]"),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Dependency]")," upgrade cobra to 1.6.0"),Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Dependency]")," upgrade logrus to 1.9.0"),Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Fixed]")," Removed builtin ",Object(l.b)("inlineCode",{parentName:"li"},"--help")," flag for subcommands. Now using ",Object(l.b)("inlineCode",{parentName:"li"},"--help")," will pas this flag to underlying ",Object(l.b)("inlineCode",{parentName:"li"},"cmd")," script.")),Object(l.b)("h2",{id:"0049"},"[0.0.49]"),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Added]")," remote mixins ",Object(l.b)("inlineCode",{parentName:"li"},"experimental")," support. See ",Object(l.b)("a",{parentName:"li",href:"/docs/config#remote-mixins-experimental"},"config")," for more details.")),Object(l.b)("h2",{id:"0048"},"[0.0.48]"),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Added]")," ",Object(l.b)("inlineCode",{parentName:"p"},"--no-depends")," global option. Lets will skip ",Object(l.b)("inlineCode",{parentName:"p"},"depends")," for running command"),Object(l.b)("pre",{parentName:"li"},Object(l.b)("code",{parentName:"pre",className:"language-shell"},"lets --no-depends run\n")))),Object(l.b)("h2",{id:"0047"},"[0.0.47]"),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Added]")," completion for command options"),Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Dependency]")," use fork of docopt.go with extended options parser")),Object(l.b)("h2",{id:"0045"},"[0.0.45]"),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Fixed]")," ",Object(l.b)("strong",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"strong"},"Breaking change"))," Fix duplicate files for checksum.\nThis will change checksum output if the same file has been read multiple times."),Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Fixed]")," Fix parsing for ref args when declared as string."),Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Added]")," ref ",Object(l.b)("inlineCode",{parentName:"li"},"args")," can be a list of string")),Object(l.b)("h2",{id:"0044"},Object(l.b)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.44"},"0.0.44")),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Fixed]")," Run ref declared in ",Object(l.b)("inlineCode",{parentName:"li"},"depends")," directive.")),Object(l.b)("h2",{id:"0043"},Object(l.b)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.43"},"0.0.43")),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Noop]")," Same as 0.0.42, deployed by accident.")),Object(l.b)("h2",{id:"0042"},Object(l.b)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.42"},"0.0.42")),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Fixed]")," Fixed publish to ",Object(l.b)("inlineCode",{parentName:"li"},"aur")," repository.")),Object(l.b)("h2",{id:"0041"},Object(l.b)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.41"},"0.0.41")),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Fixed]")," Tried to fixe publish to ",Object(l.b)("inlineCode",{parentName:"li"},"aur")," repository.")),Object(l.b)("h2",{id:"0040"},Object(l.b)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.40"},"0.0.40")),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Added]")," Allow override command arguments and env when using command in ",Object(l.b)("inlineCode",{parentName:"p"},"depends")),Object(l.b)("p",{parentName:"li"}," See example ",Object(l.b)("a",{parentName:"p",href:"/docs/config#override-arguments-in-depends-command"},"in config docs"))),Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Added]")," Validate if commands declared in ",Object(l.b)("inlineCode",{parentName:"p"},"depends")," actually exist.")),Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Refactoring]")," Refactored ",Object(l.b)("inlineCode",{parentName:"p"},"runner")," package, implemented ",Object(l.b)("inlineCode",{parentName:"p"},"Runner")," struct.")),Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Added]")," Support ",Object(l.b)("inlineCode",{parentName:"p"},"NO_COLOR")," env variable to disable colored output. See ",Object(l.b)("a",{parentName:"p",href:"https://no-color.org/"},"https://no-color.org/"))),Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Added]")," ",Object(l.b)("inlineCode",{parentName:"p"},"LETS_COMMAND_ARGS")," - will contain command's positional args. ",Object(l.b)("a",{parentName:"p",href:"/docs/env#default-environment-variables"},"See config"),"."),Object(l.b)("p",{parentName:"li"},"Also, special bash env variables such as ",Object(l.b)("inlineCode",{parentName:"p"},'"$@"')," and ",Object(l.b)("inlineCode",{parentName:"p"},'"$1"')," etc. now available inside ",Object(l.b)("inlineCode",{parentName:"p"},"cmd")," script and work as expected.")),Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Added]")," ",Object(l.b)("inlineCode",{parentName:"p"},"work_dir")," directive for command. See ",Object(l.b)("a",{parentName:"p",href:"/docs/config#work_dir"},"config"))),Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Added]")," ",Object(l.b)("inlineCode",{parentName:"p"},"shell")," directive for command. See ",Object(l.b)("a",{parentName:"p",href:"/docs/config#shell-1"},"config"))),Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Added]")," ",Object(l.b)("inlineCode",{parentName:"p"},"--init")," flag. Run ",Object(l.b)("inlineCode",{parentName:"p"},"lets --init")," to create new ",Object(l.b)("inlineCode",{parentName:"p"},"lets.yaml")," with example command")),Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Refactoring]")," updated ",Object(l.b)("inlineCode",{parentName:"p"},"bats")," test framework and adjusted all bats tests")),Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Added]")," ",Object(l.b)("inlineCode",{parentName:"p"},"ref")," directive to ",Object(l.b)("inlineCode",{parentName:"p"},"command"),". Allows to declare existing command with predefined args ",Object(l.b)("a",{parentName:"p",href:"/docs/config#ref"},"See config"),".")),Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Added]")," ",Object(l.b)("inlineCode",{parentName:"p"},"sh")," and ",Object(l.b)("inlineCode",{parentName:"p"},"checksum")," execution modes for global level ",Object(l.b)("inlineCode",{parentName:"p"},"env")," and command level ",Object(l.b)("inlineCode",{parentName:"p"},"env")," ",Object(l.b)("a",{parentName:"p",href:"/docs/config#env"},"See config"),".\n",Object(l.b)("inlineCode",{parentName:"p"},"eval_env")," is deprecated now, since ",Object(l.b)("inlineCode",{parentName:"p"},"env")," with ",Object(l.b)("inlineCode",{parentName:"p"},"sh")," execution mode does exactly the same"))),Object(l.b)("h2",{id:"0033"},Object(l.b)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.33"},"0.0.33")),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Added]")," Allow templating in command ",Object(l.b)("inlineCode",{parentName:"li"},"options")," directive ",Object(l.b)("a",{parentName:"li",href:"/docs/advanced_usage#command-templates"},"docs"))),Object(l.b)("h2",{id:"0032"},Object(l.b)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.32"},"0.0.32")),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Fixed]")," Publish lets to homebrew")),Object(l.b)("h2",{id:"0030"},Object(l.b)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.30"},"0.0.30")),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Added]")," Build ",Object(l.b)("inlineCode",{parentName:"li"},"lets")," for ",Object(l.b)("inlineCode",{parentName:"li"},"arm64 (M1)")," arch"),Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Deleted]")," Drop ",Object(l.b)("inlineCode",{parentName:"li"},"386")," arch builds"),Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Added]")," Publish ",Object(l.b)("inlineCode",{parentName:"li"},"lets")," to homebrew"),Object(l.b)("li",{parentName:"ul"},Object(l.b)("inlineCode",{parentName:"li"},"[Added]")," ",Object(l.b)("inlineCode",{parentName:"li"},"--upgrade")," flag to make self-upgrades")),Object(l.b)("h2",{id:"0029"},"0.0.29"),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Added]")," ",Object(l.b)("inlineCode",{parentName:"p"},"after")," directive to command.\nIt allows to run some script after main ",Object(l.b)("inlineCode",{parentName:"p"},"cmd")),Object(l.b)("pre",{parentName:"li"},Object(l.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  run:\n    cmd: docker-compose up redis\n    after: docker-compose stop redis\n"))),Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Added]")," ",Object(l.b)("inlineCode",{parentName:"p"},"before")," global directive to config.\nIt allows to run some script before each main ",Object(l.b)("inlineCode",{parentName:"p"},"cmd")),Object(l.b)("pre",{parentName:"li"},Object(l.b)("code",{parentName:"pre",className:"language-yaml"},"before: |\n  function @docker-compose() {\n    docker-compose --log-level ERROR $@\n  }\n\ncommands:\n  run:\n    cmd: @docker-compose up redis\n"))),Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Added]")," ignored minixs\nIt allows to include mixin only if it exists - otherwise lets will ignore it.\nUseful for git-ignored files."),Object(l.b)("p",{parentName:"li"},"Just add ",Object(l.b)("inlineCode",{parentName:"p"},"-")," prefix to mixin filename"),Object(l.b)("pre",{parentName:"li"},Object(l.b)("code",{parentName:"pre",className:"language-yaml"},"mixins:\n  - -my.yaml\n\ncommands:\n  run:\n    cmd: docker-compose up redis\n")))),Object(l.b)("h2",{id:"0028"},"0.0.28"),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Fixed]")," Added environment variable value coercion."),Object(l.b)("pre",{parentName:"li"},Object(l.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  run:\n    env:\n      VERBOSE: 1\n    cmd: docker-compose up\n")),Object(l.b)("p",{parentName:"li"},"Before 0.0.28 release this config vas invalid because ",Object(l.b)("inlineCode",{parentName:"p"},"1")," was not coerced to string ",Object(l.b)("inlineCode",{parentName:"p"},'"1"'),". Now it works as expected."))),Object(l.b)("h2",{id:"0027"},"0.0.27"),Object(l.b)("ul",null,Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},Object(l.b)("inlineCode",{parentName:"p"},"[Added]")," ",Object(l.b)("inlineCode",{parentName:"p"},"-E")," (",Object(l.b)("inlineCode",{parentName:"p"},"--env"),") command-line flag. It allows to set(override) environment variables for a running command.\nExample:"),Object(l.b)("pre",{parentName:"li"},Object(l.b)("code",{parentName:"pre",className:"language-bash"},'# lets.yaml\n...\ncommands:\n  greet:\n    env:\n      NAME: Morty\n    cmd: echo "Hello ${NAME}"\n...\n\nlets -E NAME=Rick greet\n'))),Object(l.b)("li",{parentName:"ul"},Object(l.b)("p",{parentName:"li"},"Changed behavior of ",Object(l.b)("inlineCode",{parentName:"p"},"persist_checksum")," at first run. Now, if there was no checksum and we just calculated a new checksum, that means checksum has changed, hence ",Object(l.b)("inlineCode",{parentName:"p"},"LETS_CHECKSUM_CHANGED")," will be ",Object(l.b)("inlineCode",{parentName:"p"},"true"),"."))))}p.isMDXComponent=!0}}]);