"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[80],{3905:function(e,t,n){n.d(t,{Zo:function(){return m},kt:function(){return k}});var a=n(7294);function i(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function l(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function r(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?l(Object(n),!0).forEach((function(t){i(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):l(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function p(e,t){if(null==e)return{};var n,a,i=function(e,t){if(null==e)return{};var n,a,i={},l=Object.keys(e);for(a=0;a<l.length;a++)n=l[a],t.indexOf(n)>=0||(i[n]=e[n]);return i}(e,t);if(Object.getOwnPropertySymbols){var l=Object.getOwnPropertySymbols(e);for(a=0;a<l.length;a++)n=l[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(i[n]=e[n])}return i}var o=a.createContext({}),d=function(e){var t=a.useContext(o),n=t;return e&&(n="function"==typeof e?e(t):r(r({},t),e)),n},m=function(e){var t=d(e.components);return a.createElement(o.Provider,{value:t},e.children)},s={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},c=a.forwardRef((function(e,t){var n=e.components,i=e.mdxType,l=e.originalType,o=e.parentName,m=p(e,["components","mdxType","originalType","parentName"]),c=d(n),k=i,u=c["".concat(o,".").concat(k)]||c[k]||s[k]||l;return n?a.createElement(u,r(r({ref:t},m),{},{components:n})):a.createElement(u,r({ref:t},m))}));function k(e,t){var n=arguments,i=t&&t.mdxType;if("string"==typeof e||i){var l=n.length,r=new Array(l);r[0]=c;var p={};for(var o in t)hasOwnProperty.call(t,o)&&(p[o]=t[o]);p.originalType=e,p.mdxType="string"==typeof e?e:i,r[1]=p;for(var d=2;d<l;d++)r[d]=n[d];return a.createElement.apply(null,r)}return a.createElement.apply(null,n)}c.displayName="MDXCreateElement"},1016:function(e,t,n){n.r(t),n.d(t,{assets:function(){return m},contentTitle:function(){return o},default:function(){return k},frontMatter:function(){return p},metadata:function(){return d},toc:function(){return s}});var a=n(7462),i=n(3366),l=(n(7294),n(3905)),r=["components"],p={id:"changelog",title:"Changelog"},o=void 0,d={unversionedId:"changelog",id:"changelog",title:"Changelog",description:"Unreleased",source:"@site/docs/changelog.md",sourceDirName:".",slug:"/changelog",permalink:"/docs/changelog",draft:!1,editUrl:"https://github.com/lets-cli/lets/edit/master/docs/docs/changelog.md",tags:[],version:"current",frontMatter:{id:"changelog",title:"Changelog"},sidebar:"someSidebar",previous:{title:"Best practices",permalink:"/docs/best_practices"},next:{title:"IDE/Text editors support",permalink:"/docs/ide_support"}},m={},s=[{value:"Unreleased",id:"unreleased",level:2},{value:"0.0.49",id:"0049",level:2},{value:"0.0.48",id:"0048",level:2},{value:"0.0.47",id:"0047",level:2},{value:"0.0.45",id:"0045",level:2},{value:"0.0.44",id:"0044",level:2},{value:"0.0.43",id:"0043",level:2},{value:"0.0.42",id:"0042",level:2},{value:"0.0.41",id:"0041",level:2},{value:"0.0.40",id:"0040",level:2},{value:"0.0.33",id:"0033",level:2},{value:"0.0.32",id:"0032",level:2},{value:"0.0.30",id:"0030",level:2},{value:"0.0.29",id:"0029",level:2},{value:"0.0.28",id:"0028",level:2},{value:"0.0.27",id:"0027",level:2}],c={toc:s};function k(e){var t=e.components,n=(0,i.Z)(e,r);return(0,l.kt)("wrapper",(0,a.Z)({},c,n,{components:t,mdxType:"MDXLayout"}),(0,l.kt)("h2",{id:"unreleased"},(0,l.kt)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.X"},"Unreleased")),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Dependency]")," upgrade cobra to 1.6.0")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Dependency]")," upgrade logrus to 1.9.0")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Fixed]")," Removed builtin ",(0,l.kt)("inlineCode",{parentName:"p"},"--help")," flag for subcommands. Now using ",(0,l.kt)("inlineCode",{parentName:"p"},"--help")," will pass this flag to underlying ",(0,l.kt)("inlineCode",{parentName:"p"},"cmd")," script.")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," Add ",(0,l.kt)("inlineCode",{parentName:"p"},"--debug")," (",(0,l.kt)("inlineCode",{parentName:"p"},"-d"),") debug flag. It works same as ",(0,l.kt)("inlineCode",{parentName:"p"},"LETS_DEBUG=1")," env variable. It can be specified as ",(0,l.kt)("inlineCode",{parentName:"p"},"-dd")," (or ",(0,l.kt)("inlineCode",{parentName:"p"},"LETS_DEBUG=2"),"). Lets then prints more verbose logs.")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," Add ",(0,l.kt)("inlineCode",{parentName:"p"},"--config")," ",(0,l.kt)("inlineCode",{parentName:"p"},"-c")," flag. It works same as ",(0,l.kt)("inlineCode",{parentName:"p"},"LETS_CONFIG=<path to lets file>")," env variable.")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," ",(0,l.kt)("inlineCode",{parentName:"p"},"LETS_CONFIG")," env variable now present at command runtime, and contains lets config filename. Default is ",(0,l.kt)("inlineCode",{parentName:"p"},"lets.yaml"),".")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," ",(0,l.kt)("inlineCode",{parentName:"p"},"LETS_CONFIG_DIR")," env variable now present at command runtime, and contains absolute path to dir where lets config found.")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," ",(0,l.kt)("inlineCode",{parentName:"p"},"LETS_COMMAND_WORKDIR")," env variable now present at command runtime, and contains absolute path to dir where ",(0,l.kt)("inlineCode",{parentName:"p"},"command.work_dir")," points.")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," Add ",(0,l.kt)("inlineCode",{parentName:"p"},"init")," directive to config. It is a script that will be executed only once before any other commands. It differs from ",(0,l.kt)("inlineCode",{parentName:"p"},"before")," in a way that ",(0,l.kt)("inlineCode",{parentName:"p"},"before")," is a script that is prepended to each command's script and thus will be execured every time a command executes.")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Refactoring]")," Config parsing is reimplemented using ",(0,l.kt)("inlineCode",{parentName:"p"},"UnmarhallYAML"),". This ends up in reduced size and complexity of parsing code.")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Refactoring]")," ",(0,l.kt)("inlineCode",{parentName:"p"},"Command")," now is clonable and this opened a possibility to reimplement ",(0,l.kt)("inlineCode",{parentName:"p"},"ref"),", ",(0,l.kt)("inlineCode",{parentName:"p"},"depends")," as map and ",(0,l.kt)("inlineCode",{parentName:"p"},"--no-depends")," - now we clone a command and modify a brand new struct instead of mutating the same command (which was not safe).")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Refactoring]")," ",(0,l.kt)("inlineCode",{parentName:"p"},"Command.Cmd")," script was replaced with ",(0,l.kt)("inlineCode",{parentName:"p"},"Cmds")," struct which represents a list of ",(0,l.kt)("inlineCode",{parentName:"p"},"Cmd"),". This allowed generalizing so-called cmd-as-map into a list of commands that will be executed in parallel (see ",(0,l.kt)("inlineCode",{parentName:"p"},"Executor.executeParallel"),").")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Refactoring]")," Error reporting has changed in some places and if one is depending on particular error messages it probably will break.")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Refactoring]")," Simplified ",(0,l.kt)("inlineCode",{parentName:"p"},"Executor")," by extracting commands filtering by ",(0,l.kt)("inlineCode",{parentName:"p"},"--only")," and ",(0,l.kt)("inlineCode",{parentName:"p"},"--exclude")," flags into ",(0,l.kt)("inlineCode",{parentName:"p"},"subcommand.go"),".")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," Command short syntax. See ",(0,l.kt)("a",{parentName:"p",href:"/docs/config#short-syntax"},"config reference for short syntax"),". Example:"),(0,l.kt)("p",{parentName:"li"},"Before:"),(0,l.kt)("pre",{parentName:"li"},(0,l.kt)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  hello:\n    cmd: echo Hello\n")),(0,l.kt)("p",{parentName:"li"},"After:"),(0,l.kt)("pre",{parentName:"li"},(0,l.kt)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  hello: echo Hello\n")))),(0,l.kt)("h2",{id:"0049"},(0,l.kt)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.49"},"0.0.49")),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Added]")," remote mixins ",(0,l.kt)("inlineCode",{parentName:"li"},"experimental")," support. See ",(0,l.kt)("a",{parentName:"li",href:"/docs/config#remote-mixins-experimental"},"config")," for more details.")),(0,l.kt)("h2",{id:"0048"},(0,l.kt)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.48"},"0.0.48")),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," ",(0,l.kt)("inlineCode",{parentName:"p"},"--no-depends")," global option. Lets will skip ",(0,l.kt)("inlineCode",{parentName:"p"},"depends")," for running command"),(0,l.kt)("pre",{parentName:"li"},(0,l.kt)("code",{parentName:"pre",className:"language-shell"},"lets --no-depends run\n")))),(0,l.kt)("h2",{id:"0047"},(0,l.kt)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.47"},"0.0.47")),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Added]")," completion for command options"),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Dependency]")," use fork of docopt.go with extended options parser")),(0,l.kt)("h2",{id:"0045"},(0,l.kt)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.45"},"0.0.45")),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Fixed]")," ",(0,l.kt)("strong",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"strong"},"Breaking change"))," Fix duplicate files for checksum.\nThis will change checksum output if the same file has been read multiple times."),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Fixed]")," Fix parsing for ref args when declared as string."),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Added]")," ref ",(0,l.kt)("inlineCode",{parentName:"li"},"args")," can be a list of string")),(0,l.kt)("h2",{id:"0044"},(0,l.kt)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.44"},"0.0.44")),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Fixed]")," Run ref declared in ",(0,l.kt)("inlineCode",{parentName:"li"},"depends")," directive.")),(0,l.kt)("h2",{id:"0043"},(0,l.kt)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.43"},"0.0.43")),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Noop]")," Same as 0.0.42, deployed by accident.")),(0,l.kt)("h2",{id:"0042"},(0,l.kt)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.42"},"0.0.42")),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Fixed]")," Fixed publish to ",(0,l.kt)("inlineCode",{parentName:"li"},"aur")," repository.")),(0,l.kt)("h2",{id:"0041"},(0,l.kt)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.41"},"0.0.41")),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Fixed]")," Tried to fixe publish to ",(0,l.kt)("inlineCode",{parentName:"li"},"aur")," repository.")),(0,l.kt)("h2",{id:"0040"},(0,l.kt)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.40"},"0.0.40")),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," Allow override command arguments and env when using command in ",(0,l.kt)("inlineCode",{parentName:"p"},"depends")),(0,l.kt)("p",{parentName:"li"}," See example ",(0,l.kt)("a",{parentName:"p",href:"/docs/config#override-arguments-in-depends-command"},"in config docs"))),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," Validate if commands declared in ",(0,l.kt)("inlineCode",{parentName:"p"},"depends")," actually exist.")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Refactoring]")," Refactored ",(0,l.kt)("inlineCode",{parentName:"p"},"executor")," package, implemented ",(0,l.kt)("inlineCode",{parentName:"p"},"Executor")," struct.")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," Support ",(0,l.kt)("inlineCode",{parentName:"p"},"NO_COLOR")," env variable to disable colored output. See ",(0,l.kt)("a",{parentName:"p",href:"https://no-color.org/"},"https://no-color.org/"))),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," ",(0,l.kt)("inlineCode",{parentName:"p"},"LETS_COMMAND_ARGS")," - will contain command's positional args. ",(0,l.kt)("a",{parentName:"p",href:"/docs/env#default-environment-variables"},"See config"),"."),(0,l.kt)("p",{parentName:"li"},"Also, special bash env variables such as ",(0,l.kt)("inlineCode",{parentName:"p"},'"$@"')," and ",(0,l.kt)("inlineCode",{parentName:"p"},'"$1"')," etc. now available inside ",(0,l.kt)("inlineCode",{parentName:"p"},"cmd")," script and work as expected.")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," ",(0,l.kt)("inlineCode",{parentName:"p"},"work_dir")," directive for command. See ",(0,l.kt)("a",{parentName:"p",href:"/docs/config#work_dir"},"config"))),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," ",(0,l.kt)("inlineCode",{parentName:"p"},"shell")," directive for command. See ",(0,l.kt)("a",{parentName:"p",href:"/docs/config#shell-1"},"config"))),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," ",(0,l.kt)("inlineCode",{parentName:"p"},"--init")," flag. Run ",(0,l.kt)("inlineCode",{parentName:"p"},"lets --init")," to create new ",(0,l.kt)("inlineCode",{parentName:"p"},"lets.yaml")," with example command")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Refactoring]")," updated ",(0,l.kt)("inlineCode",{parentName:"p"},"bats")," test framework and adjusted all bats tests")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," ",(0,l.kt)("inlineCode",{parentName:"p"},"ref")," directive to ",(0,l.kt)("inlineCode",{parentName:"p"},"command"),". Allows to declare existing command with predefined args ",(0,l.kt)("a",{parentName:"p",href:"/docs/config#ref"},"See config"),".")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," ",(0,l.kt)("inlineCode",{parentName:"p"},"sh")," and ",(0,l.kt)("inlineCode",{parentName:"p"},"checksum")," execution modes for global level ",(0,l.kt)("inlineCode",{parentName:"p"},"env")," and command level ",(0,l.kt)("inlineCode",{parentName:"p"},"env")," ",(0,l.kt)("a",{parentName:"p",href:"/docs/config#env"},"See config"),".\n",(0,l.kt)("inlineCode",{parentName:"p"},"eval_env")," is deprecated now, since ",(0,l.kt)("inlineCode",{parentName:"p"},"env")," with ",(0,l.kt)("inlineCode",{parentName:"p"},"sh")," execution mode does exactly the same"))),(0,l.kt)("h2",{id:"0033"},(0,l.kt)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.33"},"0.0.33")),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Added]")," Allow templating in command ",(0,l.kt)("inlineCode",{parentName:"li"},"options")," directive ",(0,l.kt)("a",{parentName:"li",href:"/docs/advanced_usage#command-templates"},"docs"))),(0,l.kt)("h2",{id:"0032"},(0,l.kt)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.32"},"0.0.32")),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Fixed]")," Publish lets to homebrew")),(0,l.kt)("h2",{id:"0030"},(0,l.kt)("a",{parentName:"h2",href:"https://github.com/lets-cli/lets/releases/tag/v0.0.30"},"0.0.30")),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Added]")," Build ",(0,l.kt)("inlineCode",{parentName:"li"},"lets")," for ",(0,l.kt)("inlineCode",{parentName:"li"},"arm64 (M1)")," arch"),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Deleted]")," Drop ",(0,l.kt)("inlineCode",{parentName:"li"},"386")," arch builds"),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Added]")," Publish ",(0,l.kt)("inlineCode",{parentName:"li"},"lets")," to homebrew"),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("inlineCode",{parentName:"li"},"[Added]")," ",(0,l.kt)("inlineCode",{parentName:"li"},"--upgrade")," flag to make self-upgrades")),(0,l.kt)("h2",{id:"0029"},"0.0.29"),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," ",(0,l.kt)("inlineCode",{parentName:"p"},"after")," directive to command.\nIt allows to run some script after main ",(0,l.kt)("inlineCode",{parentName:"p"},"cmd")),(0,l.kt)("pre",{parentName:"li"},(0,l.kt)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  run:\n    cmd: docker-compose up redis\n    after: docker-compose stop redis\n"))),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," ",(0,l.kt)("inlineCode",{parentName:"p"},"before")," global directive to config.\nIt allows to run some script before each main ",(0,l.kt)("inlineCode",{parentName:"p"},"cmd")),(0,l.kt)("pre",{parentName:"li"},(0,l.kt)("code",{parentName:"pre",className:"language-yaml"},"before: |\n  function @docker-compose() {\n    docker-compose --log-level ERROR $@\n  }\n\ncommands:\n  run:\n    cmd: @docker-compose up redis\n"))),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," ignored minixs\nIt allows to include mixin only if it exists - otherwise lets will ignore it.\nUseful for git-ignored files."),(0,l.kt)("p",{parentName:"li"},"Just add ",(0,l.kt)("inlineCode",{parentName:"p"},"-")," prefix to mixin filename"),(0,l.kt)("pre",{parentName:"li"},(0,l.kt)("code",{parentName:"pre",className:"language-yaml"},"mixins:\n  - -my.yaml\n\ncommands:\n  run:\n    cmd: docker-compose up redis\n")))),(0,l.kt)("h2",{id:"0028"},"0.0.28"),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Fixed]")," Added environment variable value coercion."),(0,l.kt)("pre",{parentName:"li"},(0,l.kt)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  run:\n    env:\n      VERBOSE: 1\n    cmd: docker-compose up\n")),(0,l.kt)("p",{parentName:"li"},"Before 0.0.28 release this config vas invalid because ",(0,l.kt)("inlineCode",{parentName:"p"},"1")," was not coerced to string ",(0,l.kt)("inlineCode",{parentName:"p"},'"1"'),". Now it works as expected."))),(0,l.kt)("h2",{id:"0027"},"0.0.27"),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("inlineCode",{parentName:"p"},"[Added]")," ",(0,l.kt)("inlineCode",{parentName:"p"},"-E")," (",(0,l.kt)("inlineCode",{parentName:"p"},"--env"),") command-line flag. It allows to set(override) environment variables for a running command.\nExample:"),(0,l.kt)("pre",{parentName:"li"},(0,l.kt)("code",{parentName:"pre",className:"language-bash"},'# lets.yaml\n...\ncommands:\n  greet:\n    env:\n      NAME: Morty\n    cmd: echo "Hello ${NAME}"\n...\n\nlets -E NAME=Rick greet\n'))),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},"Changed behavior of ",(0,l.kt)("inlineCode",{parentName:"p"},"persist_checksum")," at first run. Now, if there was no checksum and we just calculated a new checksum, that means checksum has changed, hence ",(0,l.kt)("inlineCode",{parentName:"p"},"LETS_CHECKSUM_CHANGED")," will be ",(0,l.kt)("inlineCode",{parentName:"p"},"true"),"."))))}k.isMDXComponent=!0}}]);