# directory-monitor

Tool to monitor directories if files are added to directories. Its main purpose is to check if files are added to a backup regularly to ensure that the backup job is not broken. 

## Features

* Checks if files were added since the last run.
* Sends push notifications if something needs the user's attention.
* Will only send push notifications if configurable thresholds are surpassed and thus not spam you with notifications.
* Writes a log file.

## How-To

0. Use a scheduler (e.g. cron) to schedule the tool to run regularly.
0. Create an account at https://pushover.net/, create an application and copy your User Key (userToken) and API Key (appToken).
0. Configure the application using a config file
    ```json
    {
      "heartbeatThresholdInHours": 168, 
      "deadbeatThresholdInHours": 72,
      "pushover": {
        "appToken": "secret",
        "userToken": "evenmoresecret"
      },
      "dirs": [
        "/home/fox/foo",
        "/home/fox/bar"
      ]
    }
    ``` 
    The `heartbeatThresholdInHours` is the threshold before the directory monitor will send you a notification that everything is working fine. Its purpose is to signal that the tool is running properly.
    The `deadbeatThresholdInHours` is the threshold before the directory monitor will send you a notification that the monitored directories are still broken. 

## Program Arguments

By default the config file is supposed to be at `${HOME}/.go/`. This can be overridden by providing a program argument. In the same directory as the config file, a data store will be created. It is a simple json file. The log file will also be created there.

