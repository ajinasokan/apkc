**apkc** is a bare-bones Android app build system. It lets you build simple native Android apps in APK or AAB format without entire Android IDE. It does a subset of what Android Gradle plugin does. This makes it quick and runs with low resources. apkc is mostly suitable for learning purposes.

## Features

- Base project has zero dependencies. APK is 6 KB.
- Builds base project in 2 seconds.
- Can use external JAR libraries.
- Can use external native shared libs built with C/Go/Rust etc.
- Code completion can be setup in VSCode or other editors.

## Getting started

1. Install JDK and Android SDK and make sure to have `ANDROID_HOME` and `JAVA_HOME` env variables setup.
2. Download `apkc` from Github releases for your OS/ARCH.
3. Use `doctor` command to verify all the paths are correct.

```shell
$ apkc doctor

[doctor] java /usr/local/opt/openjdk@11/ 
[doctor] javac /usr/local/opt/openjdk@11/bin/javac 
[doctor] jarsigner /usr/local/opt/openjdk@11/bin/jarsigner 
[doctor] sdk /Users/ajinasokan/Library/Android/sdk 
[doctor] aapt /Users/ajinasokan/Library/Android/sdk/build-tools/30.0.3/aapt 
[doctor] aapt2 /Users/ajinasokan/Library/Android/sdk/build-tools/30.0.3/aapt2 
[doctor] d8 /Users/ajinasokan/Library/Android/sdk/build-tools/30.0.3/d8 
[doctor] zipalign /Users/ajinasokan/Library/Android/sdk/build-tools/30.0.3/zipalign 
[doctor] android jar /Users/ajinasokan/Library/Android/sdk/platforms/android-30/android.jar
```
4. Create a hello world project.

```shell
$ apkc create myapp

Package name (com.myapp): com.ajinasokan.myapp
App name (My App): Hello World
```

5. Build APK and run in connected device. Logs from app will be streamed after successful run.

```shell
$ cd myapp

$ apkc build

[build] compiling res/layout/main.xml 
[build] compiling res/values/styles.xml 
[build] bundling resources 
[build] compiling java files 
[build] bundling classes and jars 
[build] bundling native libs 
[build] no native libs 
[build] signing apk 
[build] running zipalign 
[apkc] finished in 1.580135483s

$ apkc run

[run] installing in device(serialabcxyz) 
[run] launching main activity 

07-06 11:23:29.894 22390 22390 I libc    : SetHeapTaggingLevel: tag level set to 0
07-06 11:23:29.933 22390 22390 E jinasokan.myap: Not starting debugger since process cannot load the jdwp agent.
....
....
```

To build AAB run `apkc build --aab`. Output file will be at `build/app.{apk/aab}`.

You can add external jar dependencies to `myapp/jar` directory and native libs to `myapp/lib` in appropriate architecture sub directories.

## Code completion in VSCode

1. Install [Java extension](https://marketplace.visualstudio.com/items?itemName=redhat.java)
2. Add `android.jar` path from `doctor` command output to `referencedLibraries` list of your workspace VSCode settings file. `myapp/.vscode/settings.json`

```json
{
    "java.project.referencedLibraries": [
        "<pathtosdk>/platforms/android-30/android.jar"
    ]
}
```
3. You could also add other external jar file dependencies

## Backstory

apkc was born in 2016 as a bash script for some of [my experiments](https://ajinasokan.com/posts/smallest-app/). It then evolved into a Makefile and was used for production builds of one of [my apps](https://play.google.com/store/apps/details?id=com.innoventionist.artham). Since the Play Store started enforcing AAB files for publishing I had to move away from apkc. This project is a rewrite of the Makefile, without the incremental compilation. Incremental compilation is in the todo list.