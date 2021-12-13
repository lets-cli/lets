(window.webpackJsonp=window.webpackJsonp||[]).push([[13],{106:function(e,n,t){"use strict";t.d(n,"a",(function(){return p})),t.d(n,"b",(function(){return u}));var a=t(0),i=t.n(a);function r(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function o(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);n&&(a=a.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,a)}return t}function c(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?o(Object(t),!0).forEach((function(n){r(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):o(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function l(e,n){if(null==e)return{};var t,a,i=function(e,n){if(null==e)return{};var t,a,i={},r=Object.keys(e);for(a=0;a<r.length;a++)t=r[a],n.indexOf(t)>=0||(i[t]=e[t]);return i}(e,n);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);for(a=0;a<r.length;a++)t=r[a],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(i[t]=e[t])}return i}var b=i.a.createContext({}),s=function(e){var n=i.a.useContext(b),t=n;return e&&(t="function"==typeof e?e(n):c(c({},n),e)),t},p=function(e){var n=s(e.components);return i.a.createElement(b.Provider,{value:n},e.children)},m={inlineCode:"code",wrapper:function(e){var n=e.children;return i.a.createElement(i.a.Fragment,{},n)}},d=i.a.forwardRef((function(e,n){var t=e.components,a=e.mdxType,r=e.originalType,o=e.parentName,b=l(e,["components","mdxType","originalType","parentName"]),p=s(t),d=a,u=p["".concat(o,".").concat(d)]||p[d]||m[d]||r;return t?i.a.createElement(u,c(c({ref:n},b),{},{components:t})):i.a.createElement(u,c({ref:n},b))}));function u(e,n){var t=arguments,a=n&&n.mdxType;if("string"==typeof e||a){var r=t.length,o=new Array(r);o[0]=d;var c={};for(var l in n)hasOwnProperty.call(n,l)&&(c[l]=n[l]);c.originalType=e,c.mdxType="string"==typeof e?e:a,o[1]=c;for(var b=2;b<r;b++)o[b]=t[b];return i.a.createElement.apply(null,o)}return i.a.createElement.apply(null,t)}d.displayName="MDXCreateElement"},153:function(e,n,t){"use strict";t.r(n),n.default=t.p+"assets/images/lets-architecture-diagram-ac44548f80a96c907e3331fe90c31144.png"},83:function(e,n,t){"use strict";t.r(n),t.d(n,"frontMatter",(function(){return o})),t.d(n,"metadata",(function(){return c})),t.d(n,"toc",(function(){return l})),t.d(n,"default",(function(){return s}));var a=t(3),i=t(7),r=(t(0),t(106)),o={id:"architecture",title:"Architecture"},c={unversionedId:"architecture",id:"architecture",isDocsHomePage:!1,title:"Architecture",description:"Architecture diagram",source:"@site/docs/architecture.md",slug:"/architecture",permalink:"/docs/architecture",editUrl:"https://github.com/lets-cli/lets/edit/master/docs/docs/architecture.md",version:"current",sidebar:"someSidebar",previous:{title:"IDE/Text editors support",permalink:"/docs/ide_support"},next:{title:"Development",permalink:"/docs/development"}},l=[{value:"Parser",id:"parser",children:[{value:"How parsing works ?",id:"how-parsing-works-",children:[]},{value:"Validation",id:"validation",children:[]}]},{value:"Cobra CLI Framework",id:"cobra-cli-framework",children:[{value:"Binding our config with Cobra",id:"binding-our-config-with-cobra",children:[]}]},{value:"Runner",id:"runner",children:[]}],b={toc:l};function s(e){var n=e.components,o=Object(i.a)(e,["components"]);return Object(r.b)("wrapper",Object(a.a)({},b,o,{components:n,mdxType:"MDXLayout"}),Object(r.b)("p",null,Object(r.b)("img",{alt:"Architecture diagram",src:t(153).default})),Object(r.b)("h2",{id:"parser"},"Parser"),Object(r.b)("p",null,"At the start of lets application, parser tries to find ",Object(r.b)("inlineCode",{parentName:"p"},"lets.yaml")," file starting from current directory up to the ",Object(r.b)("inlineCode",{parentName:"p"},"/"),"."),Object(r.b)("p",null,"When config file is found, parser tries to read/parse and validate yaml config."),Object(r.b)("h4",{id:"mixins"},"Mixins"),Object(r.b)("p",null,"Lets has feature called ",Object(r.b)("a",{parentName:"p",href:"/docs/config#mixins"},"mixins"),". When parser meets ",Object(r.b)("inlineCode",{parentName:"p"},"mixins")," directive,\nit basically repeats all read/parse logic on minix files."),Object(r.b)("p",null,"Since mixin config files have some limitations, although they are parsed the same way, validation is a bit different."),Object(r.b)("h3",{id:"how-parsing-works-"},"How parsing works ?"),Object(r.b)("p",null,Object(r.b)("inlineCode",{parentName:"p"},"config.go:Config")," struct implements ",Object(r.b)("inlineCode",{parentName:"p"},"UnmarshalYAML")," function, so when ",Object(r.b)("inlineCode",{parentName:"p"},"yaml.Unmarshal")," called with ",Object(r.b)("inlineCode",{parentName:"p"},"Config")," instance passed in,\ncustom unmarshalling code is executed."),Object(r.b)("p",null,"Its common to make some normalization of commands and its data during parsing phase so the rest of the code\ndoes not have to do any kind of normalization on its own."),Object(r.b)("h3",{id:"validation"},"Validation"),Object(r.b)("p",null,"There are two validation phases."),Object(r.b)("p",null,"First validation phase happens during unmarshalling and checks if:"),Object(r.b)("ul",null,Object(r.b)("li",{parentName:"ul"},"directives names valid"),Object(r.b)("li",{parentName:"ul"},"directives types valid (array, map, string, number, etc.)"),Object(r.b)("li",{parentName:"ul"},"references to command in ",Object(r.b)("inlineCode",{parentName:"li"},"depends")," directive points to existing commands ")),Object(r.b)("p",null,"Second phase happens after we ensured that config is syntactically and semantically correct."),Object(r.b)("p",null,"Int the second phase we are checking:"),Object(r.b)("ul",null,Object(r.b)("li",{parentName:"ul"},"config version"),Object(r.b)("li",{parentName:"ul"},"circular dependencies in commands")),Object(r.b)("h2",{id:"cobra-cli-framework"},"Cobra CLI Framework"),Object(r.b)("p",null,"We are using ",Object(r.b)("inlineCode",{parentName:"p"},"Cobra")," CLI framework and delegating to it most of the work related to parsing\ncommand line arguments, help messages etc."),Object(r.b)("h3",{id:"binding-our-config-with-cobra"},"Binding our config with Cobra"),Object(r.b)("p",null,"Now we have to bind our config to ",Object(r.b)("inlineCode",{parentName:"p"},"Cobra"),"."),Object(r.b)("p",null,"Cobra has a concept of ",Object(r.b)("inlineCode",{parentName:"p"},"cobra.Command"),". It is a representation of command in CLI application, for example:"),Object(r.b)("pre",null,Object(r.b)("code",{parentName:"pre",className:"language-shell"},"git commit\ngit pull\n")),Object(r.b)("p",null,Object(r.b)("inlineCode",{parentName:"p"},"git")," is a CLI applications and\n",Object(r.b)("inlineCode",{parentName:"p"},"commit")," and ",Object(r.b)("inlineCode",{parentName:"p"},"pull")," are commands."),Object(r.b)("p",null,"In a traditional ",Object(r.b)("inlineCode",{parentName:"p"},"lets")," application commands will be what is declared in ",Object(r.b)("inlineCode",{parentName:"p"},"lets.yaml")," commands section."),Object(r.b)("p",null,"To achieve this we are creating so-called ",Object(r.b)("inlineCode",{parentName:"p"},"root")," command and ",Object(r.b)("inlineCode",{parentName:"p"},"subcommands")," from config."),Object(r.b)("h4",{id:"root-command"},"Root command"),Object(r.b)("p",null,"Root command is responsible for:"),Object(r.b)("ul",null,Object(r.b)("li",{parentName:"ul"},Object(r.b)("inlineCode",{parentName:"li"},"lets")," own command line flags such as ",Object(r.b)("inlineCode",{parentName:"li"},"--version"),", ",Object(r.b)("inlineCode",{parentName:"li"},"--upgrade"),", ",Object(r.b)("inlineCode",{parentName:"li"},"--help")," and so on."),Object(r.b)("li",{parentName:"ul"},Object(r.b)("inlineCode",{parentName:"li"},"lets")," commands autocompletion in terminal")),Object(r.b)("h4",{id:"subcommands"},"Subcommands"),Object(r.b)("p",null,"Subcommand is created from our ",Object(r.b)("inlineCode",{parentName:"p"},"Config.Commands")," (see ",Object(r.b)("inlineCode",{parentName:"p"},"initSubCommands")," function)."),Object(r.b)("p",null,"In subcommand's ",Object(r.b)("inlineCode",{parentName:"p"},"RunE")," callback we are parsing/validation/normalizing command line arguments for this subcommand\nand then finally executing command with ",Object(r.b)("inlineCode",{parentName:"p"},"Runner"),"."),Object(r.b)("p",null,"Since we are using ",Object(r.b)("inlineCode",{parentName:"p"},"docopt")," as an argument parser for subcommands, we don't let ",Object(r.b)("inlineCode",{parentName:"p"},"Cobra")," parse and interpret args,\nand instead we are passing raw arguments as is to ",Object(r.b)("inlineCode",{parentName:"p"},"Runner"),"."),Object(r.b)("h2",{id:"runner"},"Runner"),Object(r.b)("p",null,Object(r.b)("inlineCode",{parentName:"p"},"Runner")," is responsible for:"),Object(r.b)("ul",null,Object(r.b)("li",{parentName:"ul"},"parsing and preparing args using ",Object(r.b)("inlineCode",{parentName:"li"},"docopt")),Object(r.b)("li",{parentName:"ul"},"calculating and storing command's checksums"),Object(r.b)("li",{parentName:"ul"},"executing other commands from ",Object(r.b)("inlineCode",{parentName:"li"},"depends")," section"),Object(r.b)("li",{parentName:"ul"},"preparing environment "),Object(r.b)("li",{parentName:"ul"},"running command in OS using ",Object(r.b)("inlineCode",{parentName:"li"},"exec.Command"))))}s.isMDXComponent=!0}}]);