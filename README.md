# messenger-bot
A messenger bot written in go

Trying to build a messenger bot.

It's written in go and deployed to Heroku.

Try it.

Clone it.

heroku create

git push heroku master

After that you will need to add variables:

heroku config:set VALIDATION_TOKEN=********

heroku config:set FB_PAGE_ACCESS_TOKEN=*********

heroku config:set APP_SECRET=*********

#References:

https://developers.facebook.com/docs/messenger-platform/quickstart

https://github.com/jw84/messenger-bot-tutorial

https://github.com/maciekmm/messenger-platform-go-sdk

#Thread Settings

https://developers.facebook.com/docs/messenger-platform/thread-settings

curl -X POST -H "Content-Type: application/json" -d '{
quote>   "setting_type":"greeting",
quote>   "greeting":{
quote>     "text":"Welcome to My Company!"
quote>   }
quote> }' "https://graph.facebook.com/v2.6/me/thread_settings?access_token=FB_PAGE_ACCESS_TOKEN"

{"result":"Successfully updated greeting"}


curl -X POST -H "Content-Type: application/json" -d '{
quote>   "setting_type":"call_to_actions",
quote>   "thread_state":"new_thread",
quote>   "call_to_actions":[
quote>     {
quote>       "payload":"USER_DEFINED_PAYLOAD"
quote>     }
quote>   ]
quote> }' "https://graph.facebook.com/v2.6/me/thread_settings?access_token=FB_PAGE_ACCESS_TOKEN"

{"result":"Successfully added new_thread's CTAs"}


curl -X POST -H "Content-Type: application/json" -d '{
quote>   "setting_type" : "call_to_actions",
quote>   "thread_state" : "existing_thread",
quote>   "call_to_actions":[
quote>     {
quote>       "type":"postback",
quote>       "title":"Help",
quote>       "payload":"DEVELOPER_DEFINED_PAYLOAD_FOR_HELP"
quote>     },
quote>     {
quote>       "type":"postback",
quote>       "title":"Start a New Order",
quote>       "payload":"DEVELOPER_DEFINED_PAYLOAD_FOR_START_ORDER"
quote>     },
quote>     {
quote>       "type":"web_url",
quote>       "title":"View Website",
quote>       "url":"http://example.com/"
quote>     }
quote>   ]
quote> }' "https://graph.facebook.com/v2.6/me/thread_settings?access_token=FB_PAGE_ACCESS_TOKEN"

{"result":"Successfully added structured menu CTAs"}
