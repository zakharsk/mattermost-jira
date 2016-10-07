# JIRA Mattermost Webhook Connector
 
## For Heroku administrator
You need Heroku application name in `https://<heroku_app_name>.herokuapp.com` format.

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/unco-games/mattermost-jira)

[![Build Status](https://travis-ci.org/unco-games/mattermost-jira.svg?branch=master)](https://travis-ci.org/unco-games/mattermost-jira)

## For Mattermost administrator
You need incoming webhook URL in `http://<mattermost_server>/hooks/<web_hook_id>` format. Can be copied from Mattermost config.

 - System console
 - INTEGRATIONS → Custom Integrations
    - Enable Incoming Webhooks: true
    - Enable integrations to override usernames: true
    - Enable integrations to override profile picture icons: true
    
 - Team menu (3 dots near the Team name in top-left corner at the team-screen)
 - Integrations → Incoming Webhooks
 - Add Incoming Webhook
 
## For JIRA administrator
 - JIRA Aadministration → System
 - ADVANCED → WebHooks
 - Create a WebHook:
    - URL:  https://_**heroku_app_name**_.herokuapp.com?mattermost_hook_url=_**mattermost_hook_url**_
    - Issue:
        - created: true
        - updated: true
        - deleted: true

## Other info
        
#### Tested with:
 - JIRA v6.4.8
 - Mattermost v3.4.0
 
#### ToDo:
 - Handle errors (muted now)