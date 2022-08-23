import os
import bugsnag
import requests
import json
from flask import Flask, Response, request
from flask_pydantic import validate

from models import SetWebhook, DeleteWebhook


app = Flask(__name__)
router_url = os.getenv('THIS_SERVER_HTTPS_ADDRESS')
bugsnag.configure(api_key=os.environ.get('BUGSNAG'), project_root="/",)


@app.route("/", methods=['GET'])
def index():
    return "Webhook router 1.0"


# universal telegram webhook keeper
@app.route('/bot', methods=['POST'])
def router():
    # extract GET params from request
    get_params_raw = request.args.to_dict()['params'].split(',')
    get_params = {}
    for param in get_params_raw:
        key, value = param.split('=')
        get_params[key] = value

    if not {'route_ip', 'route_port'} <= get_params.keys():
        response = {"status": "ok", "result": {"error": "required GET params are missing"}}
        return Response(status=400, content_type='application/json', response=json.dumps(response))

    # build new link for route (other GET params redirected too)
    link = f"http://{get_params['route_ip']}:{get_params['route_port']}"
    i = 0
    for key, value in get_params.items():
        if key not in ['route_ip', 'route_port']:
            link += f"/?{key}={value}" if i == 0 else f"&{key}={value}"
            i += 1

    # make route
    try:
        requests.post(url=link, json=request.get_json())
    except Exception as e:
        bugsnag.notify(e)

    return Response(status=200)


# universal telegram webhook setter
@app.route('/setWebhook', methods=['POST'])
@validate(body=SetWebhook)
def webhook_register():
    post_data = request.get_json()
    extra_params = ""
    if 'extra_params' in post_data:
        for param, value in post_data['extra_params'].items():
            extra_params += f',{param}={value}'

    link = f"https://api.telegram.org/bot{post_data['telegram_token']}/setWebhook"\
           f"?max_connections={post_data['max_connections']}"\
           f"&drop_pending_updates={post_data['drop_pending_updates']}"\
           f"&url={router_url}/bot?params=route_ip={post_data['route_ip']},route_port={post_data['route_port']}{extra_params}"
    result = requests.get(link)

    return Response(status=200, content_type='application/json', response=result)


# universal telegram webhook deleter
@app.route('/deleteWebhook', methods=['POST'])
@validate(body=DeleteWebhook)
def webhook_unregister():
    params = request.get_json()
    result = requests.get(f"https://api.telegram.org/bot{params['telegram_token']}/deleteWebhook")

    return Response(status=200, content_type='application/json', response=result)

