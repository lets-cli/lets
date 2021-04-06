(window.webpackJsonp=window.webpackJsonp||[]).push([[6],{105:function(e,n,t){"use strict";t.d(n,"a",(function(){return b})),t.d(n,"b",(function(){return m}));var a=t(0),l=t.n(a);function r(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function o(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);n&&(a=a.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,a)}return t}function c(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?o(Object(t),!0).forEach((function(n){r(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):o(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function i(e,n){if(null==e)return{};var t,a,l=function(e,n){if(null==e)return{};var t,a,l={},r=Object.keys(e);for(a=0;a<r.length;a++)t=r[a],n.indexOf(t)>=0||(l[t]=e[t]);return l}(e,n);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);for(a=0;a<r.length;a++)t=r[a],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(l[t]=e[t])}return l}var p=l.a.createContext({}),s=function(e){var n=l.a.useContext(p),t=n;return e&&(t="function"==typeof e?e(n):c(c({},n),e)),t},b=function(e){var n=s(e.components);return l.a.createElement(p.Provider,{value:n},e.children)},d={inlineCode:"code",wrapper:function(e){var n=e.children;return l.a.createElement(l.a.Fragment,{},n)}},u=l.a.forwardRef((function(e,n){var t=e.components,a=e.mdxType,r=e.originalType,o=e.parentName,p=i(e,["components","mdxType","originalType","parentName"]),b=s(t),u=a,m=b["".concat(o,".").concat(u)]||b[u]||d[u]||r;return t?l.a.createElement(m,c(c({ref:n},p),{},{components:t})):l.a.createElement(m,c({ref:n},p))}));function m(e,n){var t=arguments,a=n&&n.mdxType;if("string"==typeof e||a){var r=t.length,o=new Array(r);o[0]=u;var c={};for(var i in n)hasOwnProperty.call(n,i)&&(c[i]=n[i]);c.originalType=e,c.mdxType="string"==typeof e?e:a,o[1]=c;for(var p=2;p<r;p++)o[p]=t[p];return l.a.createElement.apply(null,o)}return l.a.createElement.apply(null,t)}u.displayName="MDXCreateElement"},76:function(e,n,t){"use strict";t.r(n),t.d(n,"frontMatter",(function(){return o})),t.d(n,"metadata",(function(){return c})),t.d(n,"toc",(function(){return i})),t.d(n,"default",(function(){return s}));var a=t(3),l=t(7),r=(t(0),t(105)),o={id:"advanced_usage",title:"Advanced usage"},c={unversionedId:"advanced_usage",id:"advanced_usage",isDocsHomePage:!1,title:"Advanced usage",description:"In advanced usage we will start with a clean project and then we will add more commands to show how you can improve a developer experience in your project.",source:"@site/docs/advanced_usage.md",slug:"/advanced_usage",permalink:"/docs/advanced_usage",editUrl:"https://github.com/lets-cli/lets/edit/master/docs/docs/advanced_usage.md",version:"current",sidebar:"someSidebar",previous:{title:"Basic usage",permalink:"/docs/basic_usage"},next:{title:"Config reference",permalink:"/docs/config"}},i=[{value:"Env",id:"env",children:[]},{value:"Eval env",id:"eval-env",children:[]},{value:"Depends",id:"depends",children:[]},{value:"Checksum",id:"checksum",children:[]},{value:"Cmd as array",id:"cmd-as-array",children:[]},{value:"Options",id:"options",children:[]},{value:"Examples",id:"examples",children:[]}],p={toc:i};function s(e){var n=e.components,t=Object(l.a)(e,["components"]);return Object(r.b)("wrapper",Object(a.a)({},p,t,{components:n,mdxType:"MDXLayout"}),Object(r.b)("p",null,"In advanced usage we will start with a clean project and then we will add more commands to show how you can improve a developer experience in your project."),Object(r.b)("p",null,"Assume you have a ",Object(r.b)("inlineCode",{parentName:"p"},"node.js")," project with a ",Object(r.b)("inlineCode",{parentName:"p"},"run")," command in ",Object(r.b)("inlineCode",{parentName:"p"},"lets.yaml")," from  ",Object(r.b)("a",{parentName:"p",href:"/docs/basic_usage"},"Basic usage")),Object(r.b)("pre",null,Object(r.b)("code",{parentName:"pre",className:"language-yaml"},"shell: bash\n\ncommands:\n  run:\n    description: Run nodejs server\n    cmd: npm run server\n")),Object(r.b)("h3",{id:"env"},"Env"),Object(r.b)("p",null,"You can add global or per-command ",Object(r.b)("inlineCode",{parentName:"p"},"env"),":"),Object(r.b)("pre",null,Object(r.b)("code",{parentName:"pre",className:"language-yaml"},'shell: bash\n\nenv:\n  DEBUG: "0"\n\ncommands:\n  run:\n    description: Run nodejs server\n    env:\n      NODE_ENV: development\n    cmd: npm run server\n')),Object(r.b)("h3",{id:"eval-env"},"Eval env"),Object(r.b)("p",null,"Also if the value of the environment variable must be evaluated, you can add global or per-command ",Object(r.b)("inlineCode",{parentName:"p"},"eval_env"),":"),Object(r.b)("pre",null,Object(r.b)("code",{parentName:"pre",className:"language-yaml"},'shell: bash\n\nenv:\n  DEBUG: "0"\n\neval_env:\n  CURRENT_UID: echo "`id -u`:`id -g`"\n  CURRENT_USER_NAME: echo "`id -un`"\n\ncommands:\n  run:\n    description: Run nodejs server\n    env:\n      NODE_ENV: development\n    cmd: npm run server\n')),Object(r.b)("h3",{id:"depends"},"Depends"),Object(r.b)("p",null,"You already can start your application, and like any other project your's also have dependencies. Dependencies can be added or deleted to project "),Object(r.b)("p",null,"and developers have to know that there is some new dependency and it is needed to run ",Object(r.b)("inlineCode",{parentName:"p"},"npm install")," again."),Object(r.b)("p",null,"You can do this - just add a new command and make it as a ",Object(r.b)("inlineCode",{parentName:"p"},"run")," command dependency, so each time you call ",Object(r.b)("inlineCode",{parentName:"p"},"lets run")," - dependant command will execute first."),Object(r.b)("pre",null,Object(r.b)("code",{parentName:"pre",className:"language-yaml"},"shell: bash\n\ncommands:\n  build-deps:\n    description: Install project dependencies\n    cmd: npm install\n\n  run:\n    description: Run nodejs server\n    depends:\n      - build-deps\n    cmd: npm run server\n")),Object(r.b)("h3",{id:"checksum"},"Checksum"),Object(r.b)("p",null,"Now, each time you call ",Object(r.b)("inlineCode",{parentName:"p"},"lets run")," - ",Object(r.b)("inlineCode",{parentName:"p"},"build-deps")," will be executed first and this will guarantee that your dependencies are always up to date."),Object(r.b)("p",null,"But we have one downside - run ",Object(r.b)("inlineCode",{parentName:"p"},"npm install")," may take some time and we do not want to wait."),Object(r.b)("p",null,Object(r.b)("inlineCode",{parentName:"p"},"checksum")," to the rescue."),Object(r.b)("p",null,"Checksums allow you to know when some of the files have changed and made a decision based on that."),Object(r.b)("p",null,"When you add ",Object(r.b)("inlineCode",{parentName:"p"},"checksum")," directive to a command - ",Object(r.b)("inlineCode",{parentName:"p"},"lets")," will calculate checksum from all of the files listed in ",Object(r.b)("inlineCode",{parentName:"p"},"checksum")," and put ",Object(r.b)("inlineCode",{parentName:"p"},"LETS_CHECKSUM")," env variable to command env."),Object(r.b)("p",null,Object(r.b)("inlineCode",{parentName:"p"},"LETS_CHECKSUM")," will have a checksum value."),Object(r.b)("p",null,"We then can store this checksum somewhere in the file and check that stored checksum with a checksum from env."),Object(r.b)("p",null,"Fortunately, ",Object(r.b)("inlineCode",{parentName:"p"},"lets")," have an option for that - ",Object(r.b)("inlineCode",{parentName:"p"},"persist_checksum"),"."),Object(r.b)("p",null,"If ",Object(r.b)("inlineCode",{parentName:"p"},"persist_cheksum")," used with ",Object(r.b)("inlineCode",{parentName:"p"},"checksum")," ",Object(r.b)("inlineCode",{parentName:"p"},"lets")," will store new checksum to ",Object(r.b)("inlineCode",{parentName:"p"},".lets")," dir and each time you run a command ",Object(r.b)("inlineCode",{parentName:"p"},"lets")," will check if stored checksum changed from the one from env."),Object(r.b)("p",null,"While using ",Object(r.b)("inlineCode",{parentName:"p"},"persist_checksum"),", ",Object(r.b)("inlineCode",{parentName:"p"},"lets")," will add new env variable to command env - ",Object(r.b)("inlineCode",{parentName:"p"},"LETS_CHECKUM_CHANGED"),"."),Object(r.b)("p",null,"You can learn more about checksum in ",Object(r.b)("a",{parentName:"p",href:"/docs/config#checksum"},"Checksum section")),Object(r.b)("pre",null,Object(r.b)("code",{parentName:"pre",className:"language-yaml"},"shell: bash\n\ncommands:\n  build-deps:\n    description: Install project dependencies\n    checksum:\n      - package.json\n    persist_checksum: true\n    cmd: |\n      if [[ ${LETS_CHECKSUM_CHANGED} == true ]]; then\n        npm install\n      fi;\n\n  run:\n    description: Run nodejs server\n    depends:\n      - build-deps\n    cmd: npm run server\n    \n")),Object(r.b)("p",null,"So now ",Object(r.b)("inlineCode",{parentName:"p"},"npm")," install will be executed only on ",Object(r.b)("inlineCode",{parentName:"p"},"package.json")," change."),Object(r.b)("h3",{id:"cmd-as-array"},"Cmd as array"),Object(r.b)("p",null,"Now you have decided to add some frontend to your project. You decided to add a command to build js with a webpack."),Object(r.b)("p",null,Object(r.b)("inlineCode",{parentName:"p"},"lets.yaml")),Object(r.b)("pre",null,Object(r.b)("code",{parentName:"pre",className:"language-yaml"},"shell: bash\n\ncommands:\n  build-deps:\n    description: Install project dependencies\n    checksum:\n      - package.json\n    persist_checksum: true\n    cmd: |\n      if [[ ${LETS_CHECKSUM_CHANGED} == true ]]; then\n        npm install\n      fi;\n\n  run:\n    description: Run nodejs server\n    depends:\n      - build-deps\n    cmd: npm run server\n\n  \n  js:\n    description: Build project js\n    cmd: npm run static\n")),Object(r.b)("p",null,Object(r.b)("inlineCode",{parentName:"p"},"package.json")),Object(r.b)("pre",null,Object(r.b)("code",{parentName:"pre",className:"language-json"},'{\n    "scripts": {\n        "static": "webpack"\n    }\n}\n')),Object(r.b)("p",null,"Now you want to run js with some options like ",Object(r.b)("inlineCode",{parentName:"p"},"watch")," or different config."),Object(r.b)("p",null,"So lets update js command:"),Object(r.b)("pre",null,Object(r.b)("code",{parentName:"pre",className:"language-yaml"},"js:\n  description: Build project js\n  cmd: \n    - npm \n    - run \n    - static\n")),Object(r.b)("p",null,"All we made is just rewrite ",Object(r.b)("inlineCode",{parentName:"p"},"cmd")," to be an array of strings. Now all positional arguments will be appended to cmd during ",Object(r.b)("inlineCode",{parentName:"p"},"lets js")," call."),Object(r.b)("p",null,Object(r.b)("strong",{parentName:"p"},Object(r.b)("inlineCode",{parentName:"strong"},"lets js -- -w"))," - this will pass ",Object(r.b)("inlineCode",{parentName:"p"},"-w")," option to webpack in ",Object(r.b)("inlineCode",{parentName:"p"},"package.json")),Object(r.b)("h3",{id:"options"},"Options"),Object(r.b)("p",null,"Sooner or later you will come up with a convenient commands for your project."),Object(r.b)("p",null,Object(r.b)("inlineCode",{parentName:"p"},"lets")," options will help you with that."),Object(r.b)("p",null,"Now you have a couple of environments in your project. And you want to be able to run a server with different environments."),Object(r.b)("p",null,"Assume you have some configs:"),Object(r.b)("ul",null,Object(r.b)("li",{parentName:"ul"},"local.yaml"),Object(r.b)("li",{parentName:"ul"},"stg.yaml"),Object(r.b)("li",{parentName:"ul"},"prd.yaml")),Object(r.b)("p",null,"We can update ",Object(r.b)("inlineCode",{parentName:"p"},"run")," command using ",Object(r.b)("inlineCode",{parentName:"p"},"options"),":"),Object(r.b)("pre",null,Object(r.b)("code",{parentName:"pre",className:"language-yaml"},'run:\n  description: Run nodejs server\n  depends:\n    - build-deps\n  options: |\n    Usage: lets run [--stg] [--prd]  \n  cmd: |\n    CONFIG_PATH="local.yaml"\n    if [[ -n ${LETSOPT_STG} ]]; then\n        CONFIG_PATH="stg.yaml"\n    elif [[ -n ${LETSOPT_PRD} ]]; then\n        CONFIG_PATH="prd.yaml"\n    fi\n    npm run server -- config=$CONFIG_PATH\n')),Object(r.b)("p",null,Object(r.b)("inlineCode",{parentName:"p"},"options")," is a string in a ",Object(r.b)("inlineCode",{parentName:"p"},"docopt")," format - ",Object(r.b)("a",{parentName:"p",href:"http://docopt.org/"},"http://docopt.org/"),"."),Object(r.b)("p",null,Object(r.b)("inlineCode",{parentName:"p"},"lets")," knows how to parse docopt string and convert it in env variables."),Object(r.b)("p",null,"In a few words, ",Object(r.b)("inlineCode",{parentName:"p"},"lets")," will capitalize on all options, replace ",Object(r.b)("inlineCode",{parentName:"p"},"-")," with ",Object(r.b)("inlineCode",{parentName:"p"},"_"),"\nand append ",Object(r.b)("inlineCode",{parentName:"p"},"LETSOPT_")," prefix - so for ",Object(r.b)("inlineCode",{parentName:"p"},"lets run --stg")," we will get ",Object(r.b)("inlineCode",{parentName:"p"},"LETSOPT_STG")," env variable with no value as its a bool option."),Object(r.b)("p",null,"Another variant of option usage:"),Object(r.b)("pre",null,Object(r.b)("code",{parentName:"pre",className:"language-yaml"},"run:\n  description: Run nodejs server\n  depends:\n    - build-deps\n  options: |\n    Usage: lets run [--config=<config>] \n  cmd: |\n    npm run server -- config=${LETSOPT_CONFIG:-local.yaml}\n")),Object(r.b)("p",null,"In this example we also use options but unlike the previous example we using key-value options here."),Object(r.b)("p",null,"So if we call ",Object(r.b)("inlineCode",{parentName:"p"},"lets run --config stg.yaml")," - ",Object(r.b)("inlineCode",{parentName:"p"},"lets")," will create ",Object(r.b)("inlineCode",{parentName:"p"},"LETSOPT_CONFIG")," env variable with value ",Object(r.b)("strong",{parentName:"p"},"stg.yaml")),Object(r.b)("p",null,"One more example will show you another option ",Object(r.b)("inlineCode",{parentName:"p"},"LETSCLI"),"."),Object(r.b)("p",null,Object(r.b)("inlineCode",{parentName:"p"},"LETSCLI")," is just a complementary env variable ",Object(r.b)("inlineCode",{parentName:"p"},"lets")," will create for each ",Object(r.b)("inlineCode",{parentName:"p"},"LETSOPT"),". "),Object(r.b)("p",null,"So how does it works?"),Object(r.b)("p",null,"If we describe option ",Object(r.b)("inlineCode",{parentName:"p"},"Usage: lets run --stg")," we will actually get two env variables to one option:"),Object(r.b)("ul",null,Object(r.b)("li",{parentName:"ul"},Object(r.b)("inlineCode",{parentName:"li"},"LETSOPT_STG")," with no value"),Object(r.b)("li",{parentName:"ul"},Object(r.b)("inlineCode",{parentName:"li"},"LETSCLI_STG")," with value ",Object(r.b)("inlineCode",{parentName:"li"},"--stg"),". It just basically stores CLI argument as is.")),Object(r.b)("p",null,"You can learn more about options in ",Object(r.b)("a",{parentName:"p",href:"/docs/config#options"},"Options section")),Object(r.b)("h3",{id:"examples"},"Examples"),Object(r.b)("p",null,"There are a lot of variants how you can use ",Object(r.b)("inlineCode",{parentName:"p"},"lets")," in your project."),Object(r.b)("p",null,Object(r.b)("a",{parentName:"p",href:"https://github.com/lets-cli/lets/tree/master/examples"},"Here")," you will find more examples with:"),Object(r.b)("ul",null,Object(r.b)("li",{parentName:"ul"},"python"),Object(r.b)("li",{parentName:"ul"},"nodejs"),Object(r.b)("li",{parentName:"ul"},"docker")))}s.isMDXComponent=!0}}]);