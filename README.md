
# Testing framework for alias
This is a partial construction of acr-builder to be relied upon for cleaner testing of Yaml alias
additions. This will serve a support role string parsing related testing.

To use: Include inside of GOPATH\Go\src\github.com\Azure This folder should also have acr-builder
in order to provide support for their interdependence. In general this project should be essentially
analogous and save for preprocessor-test.go and test_frame.go could be directly copy pasted into 
acr-builder in order to transfer the work carried out which can the be git commited.


# Task Yaml alias additions

We have been discussing how to facilitate and improve on existing task yaml capabilities. In particular by
looking at declaring aliases within the yaml files for reusability and general convenience. In current
consideration we want something providing functionality similar to C's #define preprocessor directive,
with some added changes, in particular:

    Through a global alias OSS file
    - More succinctly express common patterns. Eg. using $commit instead of {{.Run.Commit}}
    - Provide aliases for pre-defined well known tools. Eg. using $curl instead of mcr.microsoft.com/acr-task-commands/curl:latest
        - Containers would be pulled if not present locally on execution.
        - This would also allow easier integration with acr's other acr tools like Gil's new prune feature
    For this in particular we could allow distinct global definition files for users to choose from 
    outside of the default 

    # Using local files 
    - Provide definition only files (for reusability) to be included in task yaml files in a similar
        fashion to including an import.
    
    # Using local aliases
    - Users may define new values or override global ones in the definition section.

Defined values would be accessed using the $ character as a preceding char. $$ used to escape $.
All alias could also be disabled in the define portion of the yaml or the alias identifier $
can be re defined.

Alias definition should be hierarchical, such that the last override of an alias will be the final 
considered value, that is hierarchically speaking:
    local definitions > file definitions > Global definitions

As an example we could see:

Example global definitions:
    registry : {{.Run.Registry}}
    commit : {{.Run.Commit}}
    id : {{.Run.ID}}

Example task yaml file:
```yaml
version: v1.0.0
    alias-src: # Import common definitions from external files (Read top down, further down indicates higher precedence)
        - org-defaults.yaml #Example of an organization default
        - proj-defaults.yaml #Additional definitions used by other local projects
    alias: # Local definitions (Take precedent over those from provided files) 
        - helm: cmd.azurecr.io/helm:v2.11.0-rc.2
    steps:
    # build website and func-test images, concurrently
    - build: -t $registry/hello-world:$id -f hello-world.dockerfile .
        when: ["-"]
    - build: -t hello-world-test  -f hello-world.dockerfile  .
        when: ["-"]
    # run built images to be tested
    - id: hello-world
        cmd: $registry/hello-world:$id
    - id: func-tests
        cmd: hello-world-test
        env:
        - test-target=hello-world
    # push hello-world if func-tests are successful  
    - push: 
        - $registry/hello-world:$id
    # helm deploy the updated hello-world image
    - cmd: >
        $helm update helloworld ./release/helm/ 
        --reuse-values 
        --set helloworld.image=$registry/hello-world:$id

```

Example org-defaults.yaml:
```yaml
    singularity: mcr.microsoft.com/acr-task-commands/singularity-builder:3.3
    pack: 'mcr.microsoft.com/azure-task-commands/buildpack:latest pack'
```

