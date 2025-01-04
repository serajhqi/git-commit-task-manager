On each commit:

1. Add Activity To Task
    
    commit message title and body will be considered as activity title and description.

    time log will be also included

    ```
    commit title

    commit body

    [task-213]
    [spent] 1h 22m 
    ```

    By using `~` in front of timelog the time span between the last task activity (if any) and the current will be set.  


    ```
    commit title

    commit body

    [task-213] task-status
    [spent] ~
    ```

2. Create Task

    The text in front of the code will be considered as task title and description and status. the delimiter is `~~~`.

    - if you don't include task title and description the commit title and body will be used.
    - if you don't include the status. the status will be set as `OPEN`

    Task Statuses:
    - OPEN
    - CLOSE


    to create a new task you need to include the project name.

    ```
    commit title

    commit body

    [proj-23] task title ~~~ task description ~~~ status
    [spent] 1h 22m 
    ```




2. Create Project


git commit should be checked on commit
but sending data to server should be done on push to aslo send the commit hash to server