(window.webpackJsonp=window.webpackJsonp||[]).push([[19],{107:function(e,n,t){"use strict";t.d(n,"a",(function(){return s})),t.d(n,"b",(function(){return u}));var l=t(0),a=t.n(l);function i(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function c(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var l=Object.getOwnPropertySymbols(e);n&&(l=l.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,l)}return t}function b(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?c(Object(t),!0).forEach((function(n){i(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):c(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function r(e,n){if(null==e)return{};var t,l,a=function(e,n){if(null==e)return{};var t,l,a={},i=Object.keys(e);for(l=0;l<i.length;l++)t=i[l],n.indexOf(t)>=0||(a[t]=e[t]);return a}(e,n);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(l=0;l<i.length;l++)t=i[l],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(a[t]=e[t])}return a}var p=a.a.createContext({}),o=function(e){var n=a.a.useContext(p),t=n;return e&&(t="function"==typeof e?e(n):b(b({},n),e)),t},s=function(e){var n=o(e.components);return a.a.createElement(p.Provider,{value:n},e.children)},m={inlineCode:"code",wrapper:function(e){var n=e.children;return a.a.createElement(a.a.Fragment,{},n)}},d=a.a.forwardRef((function(e,n){var t=e.components,l=e.mdxType,i=e.originalType,c=e.parentName,p=r(e,["components","mdxType","originalType","parentName"]),s=o(t),d=l,u=s["".concat(c,".").concat(d)]||s[d]||m[d]||i;return t?a.a.createElement(u,b(b({ref:n},p),{},{components:t})):a.a.createElement(u,b({ref:n},p))}));function u(e,n){var t=arguments,l=n&&n.mdxType;if("string"==typeof e||l){var i=t.length,c=new Array(i);c[0]=d;var b={};for(var r in n)hasOwnProperty.call(n,r)&&(b[r]=n[r]);b.originalType=e,b.mdxType="string"==typeof e?e:l,c[1]=b;for(var p=2;p<i;p++)c[p]=t[p];return a.a.createElement.apply(null,c)}return a.a.createElement.apply(null,t)}d.displayName="MDXCreateElement"},90:function(e,n,t){"use strict";t.r(n),t.d(n,"frontMatter",(function(){return c})),t.d(n,"metadata",(function(){return b})),t.d(n,"toc",(function(){return r})),t.d(n,"default",(function(){return o}));var l=t(3),a=t(7),i=(t(0),t(107)),c={id:"config",title:"Config reference"},b={unversionedId:"config",id:"config",isDocsHomePage:!1,title:"Config reference",description:"Config schema",source:"@site/docs/config.md",slug:"/config",permalink:"/docs/config",editUrl:"https://github.com/lets-cli/lets/edit/master/docs/docs/config.md",version:"current",sidebar:"someSidebar",previous:{title:"Advanced usage",permalink:"/docs/advanced_usage"},next:{title:"Lets cli",permalink:"/docs/cli"}},r=[{value:"Top-level directives:",id:"top-level-directives",children:[{value:"Version",id:"version",children:[]},{value:"Shell",id:"shell",children:[]},{value:"Global env",id:"global-env",children:[]},{value:"Global eval_env",id:"global-eval_env",children:[]},{value:"Global before",id:"global-before",children:[]},{value:"Mixins",id:"mixins",children:[]},{value:"Ignored mixins",id:"ignored-mixins",children:[]},{value:"Remote mixins <code>(experimental)</code>",id:"remote-mixins-experimental",children:[]},{value:"Commands",id:"commands",children:[]}]},{value:"Command directives:",id:"command-directives",children:[{value:"<code>description</code>",id:"description",children:[]},{value:"<code>cmd</code>",id:"cmd",children:[]},{value:"<code>work_dir</code>",id:"work_dir",children:[]},{value:"<code>shell</code>",id:"shell-1",children:[]},{value:"<code>after</code>",id:"after",children:[]},{value:"<code>depends</code>",id:"depends",children:[]},{value:"<code>options</code>",id:"options",children:[]},{value:"<code>env</code>",id:"env",children:[]},{value:"<code>eval_env</code>",id:"eval_env",children:[]},{value:"<code>checksum</code>",id:"checksum",children:[]},{value:"<code>persist_checksum</code>",id:"persist_checksum",children:[]},{value:"<code>ref</code>",id:"ref",children:[]},{value:"<code>args</code>",id:"args",children:[]}]}],p={toc:r};function o(e){var n=e.components,t=Object(a.a)(e,["components"]);return Object(i.b)("wrapper",Object(l.a)({},p,t,{components:n,mdxType:"MDXLayout"}),Object(i.b)("p",null,"Config schema"),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#shell"},"shell")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#mixins"},"mixins")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#global-env"},"env")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#global-eval_env"},"eval_env")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#before"},"before")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#commands"},"commands"),Object(i.b)("ul",{parentName:"li"},Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#description"},"description")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#cmd"},"cmd")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#work_dir"},"work_dir")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#after"},"after")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#depends"},"depends")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#options"},"options")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#env"},"env")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#eval_env"},"eval_env")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#checksum"},"checksum")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#persist_checksum"},"persist_checksum")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#ref"},"ref")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("a",{parentName:"li",href:"#args"},"args"))))),Object(i.b)("h2",{id:"top-level-directives"},"Top-level directives:"),Object(i.b)("h3",{id:"version"},"Version"),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: version")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: semver string")),Object(i.b)("p",null,"Specify ",Object(i.b)("strong",{parentName:"p"},"minimum required")," ",Object(i.b)("inlineCode",{parentName:"p"},"lets")," version to run this config."),Object(i.b)("p",null,"Example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"version: '0.0.20'\n")),Object(i.b)("h3",{id:"shell"},"Shell"),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: shell")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: string")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"required: true")),Object(i.b)("p",null,"Specify shell to use when running commands"),Object(i.b)("p",null,"Example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"shell: bash\n")),Object(i.b)("h3",{id:"global-env"},"Global env"),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: env")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: map string => string or map with execution mode")),Object(i.b)("p",null,"Specify global env for all commands."),Object(i.b)("p",null,"Env can be declared as static value or with execution mode:"),Object(i.b)("p",null,"Example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},'shell: bash\nenv:\n  MY_GLOBAL_ENV: "123"\n  MY_GLOBAL_ENV_2: \n    sh: echo "`id`"\n  MY_GLOBAL_ENV_3:\n    checksum: [Readme.md, package.json]\n')),Object(i.b)("h3",{id:"global-eval_env"},"Global eval_env"),Object(i.b)("p",null,Object(i.b)("strong",{parentName:"p"},Object(i.b)("inlineCode",{parentName:"strong"},"Deprecated"))),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: eval_env")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: mapping string => string")),Object(i.b)("blockquote",null,Object(i.b)("p",{parentName:"blockquote"},"Since ",Object(i.b)("inlineCode",{parentName:"p"},"env")," now has ",Object(i.b)("inlineCode",{parentName:"p"},"sh")," execution mode, ",Object(i.b)("inlineCode",{parentName:"p"},"eval_env")," is deprecated.")),Object(i.b)("p",null,"Specify global eval_env for all commands."),Object(i.b)("p",null,"Example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},'shell: bash\neval_env:\n  CURRENT_UID: echo "`id -u`:`id -g`"\n')),Object(i.b)("h3",{id:"global-before"},"Global before"),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: before")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: string")),Object(i.b)("p",null,"Specify global before script for all commands."),Object(i.b)("p",null,"Example:"),Object(i.b)("p",null,"Run ",Object(i.b)("inlineCode",{parentName:"p"},"redis")," with docker-compose using log level ERROR"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"shell: bash\n\nbefore:\n  function @docker-compose() {\n    docker-compose --log-level ERROR $@\n  }\n\n  export XXX=123\n\ncommands:\n  redis: |\n    echo $XXX\n    @docker-compose up redis\n")),Object(i.b)("h3",{id:"mixins"},"Mixins"),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: mixins")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type:")," "),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},Object(i.b)("inlineCode",{parentName:"li"},"list of strings")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("inlineCode",{parentName:"li"},"list of map"))),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"Example")),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre"},"mixins:\n  - lets.build.yaml\n  - url: https://raw.githubusercontent.com/lets-cli/lets/master/lets.build.yaml\n    version: 1\n")),Object(i.b)("p",null,"Allows to split ",Object(i.b)("inlineCode",{parentName:"p"},"lets.yaml")," into mixins (mixin config files)."),Object(i.b)("p",null,"To make ",Object(i.b)("inlineCode",{parentName:"p"},"lets.yaml")," small and readable it is convenient to split main config into many smaller ones and include them"),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"Full example")),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"# in lets.yaml\n...\nshell: bash\nmixins:\n  - test.yaml\n\ncommands:\n  echo:\n    cmd: echo Hi\n    \n# in test.yaml\n...\ncommands:\n  test:\n    cmd: echo Testing...\n")),Object(i.b)("h3",{id:"ignored-mixins"},"Ignored mixins"),Object(i.b)("p",null,"It is possible to specify mixin file which could not exist. It is convenient when you have\ngit-ignored file where you write your own commands."),Object(i.b)("p",null,"To make ",Object(i.b)("inlineCode",{parentName:"p"},"lets")," read this mixin just add ",Object(i.b)("inlineCode",{parentName:"p"},"-")," prefix to filename"),Object(i.b)("p",null,"For example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"shell: bash\nmixins:\n  - -my.yaml\n")),Object(i.b)("p",null,"Now if ",Object(i.b)("inlineCode",{parentName:"p"},"my.yaml")," exists - it will be loaded as a mixin. If it is not exist - ",Object(i.b)("inlineCode",{parentName:"p"},"lets")," will skip it."),Object(i.b)("h3",{id:"remote-mixins-experimental"},"Remote mixins ",Object(i.b)("inlineCode",{parentName:"h3"},"(experimental)")),Object(i.b)("p",null,"It is possible to specify mixin as url. Lets will download it and load it as a mixin.\nFile will be stored in ",Object(i.b)("inlineCode",{parentName:"p"},".lets/mixins")," directory."),Object(i.b)("p",null,"By default mixin filename will be sha256 hash of url."),Object(i.b)("p",null,"You can specify ",Object(i.b)("inlineCode",{parentName:"p"},"version")," key. If url is not versioned, lets will use ",Object(i.b)("inlineCode",{parentName:"p"},"version")," for filename hash as well (",Object(i.b)("inlineCode",{parentName:"p"},"hash(url) + hash(version)"),")."),Object(i.b)("p",null,"For example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"shell: bash\nmixins:\n  - url: https://raw.githubusercontent.com/lets-cli/lets/master/lets.build.yaml\n    version: 1\n")),Object(i.b)("h3",{id:"commands"},"Commands"),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: commands")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: mapping")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"required: true")),Object(i.b)("p",null,"Mapping of all available commands"),Object(i.b)("p",null,"Example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  test:\n    description: Test something\n")),Object(i.b)("h2",{id:"command-directives"},"Command directives:"),Object(i.b)("h3",{id:"description"},Object(i.b)("inlineCode",{parentName:"h3"},"description")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: description")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: string")),Object(i.b)("p",null,"Short description of command - shown in help message"),Object(i.b)("p",null,"Example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  test:\n    description: Test something\n")),Object(i.b)("h3",{id:"cmd"},Object(i.b)("inlineCode",{parentName:"h3"},"cmd")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: cmd")),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre"},"type: \n  - string\n  - array of strings\n  - map of string => string (experimental)\n")),Object(i.b)("p",null,"Actual command to run in shell."),Object(i.b)("p",null,"Can be either:"),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},"a string (also a multiline string)"),Object(i.b)("li",{parentName:"ul"},"an array of strings - it will allow to append all arguments passed to command as is (see bellow)"),Object(i.b)("li",{parentName:"ul"},"a map of string => string - this will allow run commands in parallel ",Object(i.b)("inlineCode",{parentName:"li"},"(experimental)"))),Object(i.b)("p",null,"Example single string:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  test:\n    description: Test something\n    cmd: go test ./... -v\n")),Object(i.b)("p",null,"Example multiline string:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},'commands:\n  test:\n    description: Test something\n    cmd: |\n      echo "Running go tests..."\n      go test ./... -v\n')),Object(i.b)("p",null,"Example array of strings:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  test:\n    description: Test something\n    cmd: \n      - go\n      - test\n      - ./...\n")),Object(i.b)("p",null,"When run with cmd as array of strings:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-bash"},"lets test -v\n")),Object(i.b)("p",null,"the ",Object(i.b)("inlineCode",{parentName:"p"},"-v")," will be appended, so the resulting command to run will be ",Object(i.b)("inlineCode",{parentName:"p"},"go test ./... -v")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"cmd")," can be a map ",Object(i.b)("inlineCode",{parentName:"p"},"(it is experimental feature)"),"."),Object(i.b)("p",null,"Example of map of string => string"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  run:\n    description: Test something\n    cmd: \n      app: npm run app\n      nginx: docker-compose up nginx\n      redis: docker-compsoe up redis\n")),Object(i.b)("p",null,"There are two flags ",Object(i.b)("inlineCode",{parentName:"p"},"--only")," and ",Object(i.b)("inlineCode",{parentName:"p"},"--exclude")," you can use with cmd map."),Object(i.b)("p",null,"There must be used before command name:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-bash"},"lets --only app run\n")),Object(i.b)("h3",{id:"work_dir"},Object(i.b)("inlineCode",{parentName:"h3"},"work_dir")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: work_dir")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: string")),Object(i.b)("p",null,"Specify work directory to run in. Path must be relative to project root. Be default command's workdir is project root (where lets.yaml located)."),Object(i.b)("p",null,"Example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  run-docs:\n    description: Run docusaurus documentation live\n    work_dir: docs\n    cmd: npm start\n")),Object(i.b)("h3",{id:"shell-1"},Object(i.b)("inlineCode",{parentName:"h3"},"shell")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: shell")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: string")),Object(i.b)("p",null,"Specify shell to run command in."),Object(i.b)("p",null,"Any shell can be used, not only sh-compatible, for example ",Object(i.b)("inlineCode",{parentName:"p"},"python"),"."),Object(i.b)("p",null,"Example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"shell: bash\ncommands:\n  run-sh:\n    shell: /bin/sh\n    cmd: echo Hi\n    \n  run-py:\n    shell: python\n    cmd: print('hi')\n")),Object(i.b)("h3",{id:"after"},Object(i.b)("inlineCode",{parentName:"h3"},"after")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: after")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: string")),Object(i.b)("p",null,"Specify script to run after the actual command. May be useful, when we want to cleanup some resources or stop some services"),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"after")," script is guaranteed to execute if specified, event if ",Object(i.b)("inlineCode",{parentName:"p"},"cmd")," exit code is not ",Object(i.b)("inlineCode",{parentName:"p"},"0")),Object(i.b)("p",null,"Example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  redis:\n    description: Run redis\n    cmd: docker-compose up redis\n    after: docker-compose stop redis\n\n  run:\n    description: Run app and services\n    cmd: \n      app: node server.js\n      redis: docker-compose up redis\n    after: |\n      echo Stopping app and redis\n      docker-compose stop redis\n")),Object(i.b)("h3",{id:"depends"},Object(i.b)("inlineCode",{parentName:"h3"},"depends")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: depends")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: array of string or array or object")),Object(i.b)("p",null,"Specify what commands to run before the actual command. May be useful, when you have one shared command.\nFor example, lets say you have command ",Object(i.b)("inlineCode",{parentName:"p"},"build"),", which builds docker image."),Object(i.b)("p",null,"Example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  build:\n    description: Build docker image\n    cmd: docker build -t myimg . -f Dockerfile\n\n  test:\n    description: Test something\n    depends: [build]\n    cmd: go test ./... -v\n\n  fmt:\n    description: Format the code\n    depends: [build]\n    cmd: go fmt\n")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"build")," command will be executed each time you run ",Object(i.b)("inlineCode",{parentName:"p"},"lets test")," or ",Object(i.b)("inlineCode",{parentName:"p"},"lets fmt")),Object(i.b)("h4",{id:"override-arguments-in-depends-command"},"Override arguments in depends command"),Object(i.b)("p",null,"It is possible to override arguments or env for any commands declared in depends."),Object(i.b)("p",null,"For example we want:"),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},Object(i.b)("inlineCode",{parentName:"li"},"build")," command to be executed with ",Object(i.b)("inlineCode",{parentName:"li"},"--verbose")," flag in test ",Object(i.b)("inlineCode",{parentName:"li"},"depends"),"."),Object(i.b)("li",{parentName:"ul"},Object(i.b)("inlineCode",{parentName:"li"},"alarm")," command to be executed with ",Object(i.b)("inlineCode",{parentName:"li"},"Something is happening")," arg and ",Object(i.b)("inlineCode",{parentName:"li"},"LEVEL: INFO")," env in test ",Object(i.b)("inlineCode",{parentName:"li"},"depends"),".")),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  greet:\n    cmd: echo Hi developer\n\n  alarm:\n    options: |\n      Usage: lets alarm <msg>\n    env:\n      LEVEL: DEBUG\n    cmd: echo Alarm ${LETSOPT_MSG}\n\n  build:\n    description: Build docker image\n    options: |\n      lets build [--verbose]\n    cmd: |\n      if [[ -n ${LETSOPT_VERBOSE} ]]; then\n        echo Building docker image\n      fi\n      docker build -t myimg . -f Dockerfile\n\n  test:\n    description: Test something\n    depends:\n      - greet\n      - name: alarm\n        args: Something is happening\n        env:\n          LEVEL: INFO\n      - name: build:\n        args: [--verbose]\n    cmd: go test ./... -v\n")),Object(i.b)("p",null,"Running ",Object(i.b)("inlineCode",{parentName:"p"},"lets test")," will output: "),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-shell"},"# lets test\n# Hi developer\n# Something is happening\n# Building docker image\n# ... continue building docker image\n")),Object(i.b)("h3",{id:"options"},Object(i.b)("inlineCode",{parentName:"h3"},"options")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: options")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: string (multiline string)")),Object(i.b)("p",null,"One of the most cool things about ",Object(i.b)("inlineCode",{parentName:"p"},"lets")," than it has built in docopt parsing.\nAll you need is to write a valid docopt for a command and lets will parse and inject all values for you."),Object(i.b)("p",null,"More info ",Object(i.b)("a",{parentName:"p",href:"http://docopt.org"},"http://docopt.org")),Object(i.b)("p",null,"When parsed, ",Object(i.b)("inlineCode",{parentName:"p"},"lets")," will provide two kind of env variables:"),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},Object(i.b)("inlineCode",{parentName:"li"},"LETSOPT_<VAR>")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("inlineCode",{parentName:"li"},"LETSCLI_<VAR>"))),Object(i.b)("p",null,"How does it work?"),Object(i.b)("p",null,"Lets see an example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  echo-env:\n    description: Echo env vars\n    options:\n      Usage: lets [--log-level=<level>] [--debug] <args>...\n      Options:\n        <args>...       List of required positional args\n        --log-level,-l      Log level\n        --debug,-d      Run with debug\n    cmd: |\n      echo ${LETSOPT_ARGS}\n      app ${LETSCLI_DEBUG}\n")),Object(i.b)("p",null,"So here we have:"),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"args")," - is a list of required positional args"),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"--log-level")," - is a key-value flag, must be provided with some value"),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"--debug")," - is a bool flag, if specified, means true, if no specified means false"),Object(i.b)("p",null,"In the env of ",Object(i.b)("inlineCode",{parentName:"p"},"cmd")," command there will be available two types of env variables:"),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"lets echo-env --log-level=info --debug one two three")),Object(i.b)("p",null,"Parsed and formatted key values"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-bash"},"echo LETSOPT_ARGS=${LETSOPT_ARGS} # LETSOPT_ARGS=one two three\necho LETSOPT_LOG_LEVEL=${LETSOPT_LOG_LEVEL} # LETSOPT_LOG_LEVEL=info\necho LETSOPT_DEBUG=${LETSOPT_DEBUG} # LETSOPT_DEBUG=true\n")),Object(i.b)("p",null,"Raw flags (useful if for example you want to pass ",Object(i.b)("inlineCode",{parentName:"p"},"--debug")," as is)"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-bash"},"echo LETSCLI_ARGS=${LETSCLI_ARGS} # LETSCLI_ARGS=one two three\necho LETSCLI_LOG_LEVEL=${LETSCLI_LOG_LEVEL} # LETSCLI_LOG_LEVEL=--log-level info\necho LETSCLI_DEBUG=${LETSCLI_DEBUG} # LETSCLI_DEBUG=--debug\n")),Object(i.b)("h3",{id:"env"},Object(i.b)("inlineCode",{parentName:"h3"},"env")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: env")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: mapping string => string or map with execution mode")),Object(i.b)("p",null,"Env is as simple as it sounds. Define additional env for a command: "),Object(i.b)("p",null,"Env can be declared as static value or with execution mode:"),Object(i.b)("p",null,"Example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},'commands:\n  test:\n    description: Test something\n    env:\n      GO111MODULE: "on"\n      GOPROXY: https://goproxy.io\n      MY_ENV_1:\n        sh: echo "`id`"\n      MY_ENV_2:\n        checksum: [Readme.md, package.json]\n    cmd: go build -o lets *.go\n')),Object(i.b)("h3",{id:"eval_env"},Object(i.b)("inlineCode",{parentName:"h3"},"eval_env")),Object(i.b)("p",null,Object(i.b)("strong",{parentName:"p"},Object(i.b)("inlineCode",{parentName:"strong"},"Deprecated"))),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: eval_env")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: mapping string => string")),Object(i.b)("blockquote",null,Object(i.b)("p",{parentName:"blockquote"},"Since ",Object(i.b)("inlineCode",{parentName:"p"},"env")," now has ",Object(i.b)("inlineCode",{parentName:"p"},"sh")," execution mode, ",Object(i.b)("inlineCode",{parentName:"p"},"eval_env")," is deprecated.")),Object(i.b)("p",null,"Same as env but allows you to dynamically compute env:"),Object(i.b)("p",null,"Example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},'commands:\n  test:\n    description: Test something\n    eval_env:\n      CURRENT_UID: echo "`id -u`:`id -g`"\n      CURRENT_USER_NAME: echo "`id -un`"\n    cmd: go build -o lets *.go\n')),Object(i.b)("p",null,"Value will be executed in shell and result will be saved in env."),Object(i.b)("h3",{id:"checksum"},Object(i.b)("inlineCode",{parentName:"h3"},"checksum")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: checksum")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: array of string | mapping string => array of string")),Object(i.b)("p",null,"Checksum used for computing file hashes. It is useful when you depend on some files content changes."),Object(i.b)("p",null,"In ",Object(i.b)("inlineCode",{parentName:"p"},"checksum")," you can specify:"),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},"a list of file names "),Object(i.b)("li",{parentName:"ul"},"a list of file regexp patterns (parsed via go ",Object(i.b)("inlineCode",{parentName:"li"},"path/filepath.Glob"),")")),Object(i.b)("p",null,"or"),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},"a mapping where key is name of env variable and value is:",Object(i.b)("ul",{parentName:"li"},Object(i.b)("li",{parentName:"ul"},"a list of file names "),Object(i.b)("li",{parentName:"ul"},"a list of file regexp patterns (parsed via go ",Object(i.b)("inlineCode",{parentName:"li"},"path/filepath.Glob"),")")))),Object(i.b)("p",null,"Each time a command runs, ",Object(i.b)("inlineCode",{parentName:"p"},"lets")," will calculate the checksum of all files specified in ",Object(i.b)("inlineCode",{parentName:"p"},"checksum"),"."),Object(i.b)("p",null,"Result then can be accessed via ",Object(i.b)("inlineCode",{parentName:"p"},"LETS_CHECKSUM")," env variable."),Object(i.b)("p",null,"If checksum is a mapping, e.g:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  build:\n    checksum:\n      deps:\n        - package.json\n      doc:\n        - Readme.md\n")),Object(i.b)("p",null,"Resulting env will be:"),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},Object(i.b)("inlineCode",{parentName:"li"},"LETS_CHECKSUM_DEPS")," - checksum of deps files"),Object(i.b)("li",{parentName:"ul"},Object(i.b)("inlineCode",{parentName:"li"},"LETS_CHECKSUM_DOC")," - checksum of doc files"),Object(i.b)("li",{parentName:"ul"},Object(i.b)("inlineCode",{parentName:"li"},"LETS_CHECKSUM")," - checksum of all checksums (deps and doc)")),Object(i.b)("p",null,"Checksum is calculated with ",Object(i.b)("inlineCode",{parentName:"p"},"sha1"),"."),Object(i.b)("p",null,"If you specify patterns, ",Object(i.b)("inlineCode",{parentName:"p"},"lets")," will try to find all matches and will calculate checksum of that files."),Object(i.b)("p",null,"Example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"shell: bash\ncommands:\n  app-build:\n    checksum: \n      - requirements-*.txt\n    cmd: |\n      docker pull myrepo/app:${LETS_CHECKSUM}\n      docker run --rm myrepo/app${LETS_CHECKSUM} python -m app       \n")),Object(i.b)("h3",{id:"persist_checksum"},Object(i.b)("inlineCode",{parentName:"h3"},"persist_checksum")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: persist_checksum")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: bool")),Object(i.b)("p",null,"This feature is useful when you want to know that something has changed between two executions of a command."),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"persist_checksum")," can be used only if ",Object(i.b)("inlineCode",{parentName:"p"},"checksum")," declared for command."),Object(i.b)("p",null,"If set to ",Object(i.b)("inlineCode",{parentName:"p"},"true"),", each run all calculated checksums will be stored to disk."),Object(i.b)("p",null,"After each subsequent run ",Object(i.b)("inlineCode",{parentName:"p"},"lets")," will check if new checksum and stored checksum are different."),Object(i.b)("p",null,"Result of that check will be exposed via ",Object(i.b)("inlineCode",{parentName:"p"},"LETS_CHECKSUM_CHANGED")," and ",Object(i.b)("inlineCode",{parentName:"p"},"LETS_CHECKSUM_[checksum-name]_CHANGED")," env variables. "),Object(i.b)("p",null,Object(i.b)("strong",{parentName:"p"},"IMPORTANT"),": New checksum will override old checksum only if cmd has exit code ",Object(i.b)("strong",{parentName:"p"},"0")," "),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"LETS_CHECKSUM_CHANGED")," will be true after the very first execution, because when you first run command, there is no checksum yet, so we are calculating new checksum - that means that checksum has changed."),Object(i.b)("p",null,"Example:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  build:\n    persist_checksum: true\n    checksum:\n      deps:\n        - package.json\n      doc:\n        - Readme.md\n")),Object(i.b)("p",null,"Resulting env will be:"),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},Object(i.b)("p",{parentName:"li"},Object(i.b)("inlineCode",{parentName:"p"},"LETS_CHECKSUM_DEPS")," - checksum of deps files")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("p",{parentName:"li"},Object(i.b)("inlineCode",{parentName:"p"},"LETS_CHECKSUM_DOC")," - checksum of doc files")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("p",{parentName:"li"},Object(i.b)("inlineCode",{parentName:"p"},"LETS_CHECKSUM")," - checksum of all checksums (deps and doc)")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("p",{parentName:"li"},Object(i.b)("inlineCode",{parentName:"p"},"LETS_CHECKSUM_DEPS_CHANGED")," - is checksum of deps files changed")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("p",{parentName:"li"},Object(i.b)("inlineCode",{parentName:"p"},"LETS_CHECKSUM_DOC_CHANGED")," - is checksum of doc files changed")),Object(i.b)("li",{parentName:"ul"},Object(i.b)("p",{parentName:"li"},Object(i.b)("inlineCode",{parentName:"p"},"LETS_CHECKSUM_CHANGED")," - is checksum of all checksums (deps and doc) changed"))),Object(i.b)("h3",{id:"ref"},Object(i.b)("inlineCode",{parentName:"h3"},"ref")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: ref")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: string")),Object(i.b)("p",null,Object(i.b)("strong",{parentName:"p"},Object(i.b)("inlineCode",{parentName:"strong"},"Experimental feature"))),Object(i.b)("p",null,"NOTE: ",Object(i.b)("inlineCode",{parentName:"p"},"ref")," is not compatible (yet) with any directives except ",Object(i.b)("inlineCode",{parentName:"p"},"args"),". Actually ",Object(i.b)("inlineCode",{parentName:"p"},"ref")," is a special syntax that turns command into reference to command. It may be changed in the future."),Object(i.b)("p",null,"Allows to run command with predefined arguments. Before this was implemented, if you had some commmand and wanted same command but with some predefined args, you had to use so called ",Object(i.b)("inlineCode",{parentName:"p"},"lets-in-lets")," hack."),Object(i.b)("p",null,"Before:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  ls:\n    cmd: [ls]\n\n  ls-table:\n    cmd: lets ls -l\n")),Object(i.b)("p",null,"Now:"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-yaml"},"commands:\n  hello:\n    cmd: echo Hello $@\n\n  hello-world:\n    ref: hello\n    args: World\n\n  hello-by-name:\n    ref: hello\n    args: [Dear, Friend]\n")),Object(i.b)("h3",{id:"args"},Object(i.b)("inlineCode",{parentName:"h3"},"args")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"key: args")),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"type: string or list of string")),Object(i.b)("p",null,Object(i.b)("strong",{parentName:"p"},Object(i.b)("inlineCode",{parentName:"strong"},"Experimental feature"))),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"args")," is used only with ",Object(i.b)("a",{parentName:"p",href:"#ref"},"ref")," and allows to set additional positional args to referenced command. See ",Object(i.b)("a",{parentName:"p",href:"#ref"},"ref")," example."))}o.isMDXComponent=!0}}]);