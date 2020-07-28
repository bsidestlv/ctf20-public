import json
import boto3
import os

PRISONERS = {
    1: "Next time I see you, remind me not to talk to you.",
    2: "I never forget a face, but in your case I'll be glad to make an exception.",
    3: os.environ['Flag']
}


def get_user(access_key, prisoner=None):
    client = boto3.client('cognito-idp')
    try:
        response = client.get_user(AccessToken=access_key)
    except:
        return wrap_response(299, {
            'text': 'The server has accepted your request but thinks you can Try Harder!'
        })
    for name in response['UserAttributes']:
        if 'iswarden' in name['Name'].lower():
            if name['Value'] == '0':
                return wrap_response(403, {
                    'text': 'Only the warden is allowed!'
                })
            elif name['Value'] == '1':
                if prisoner is not None:
                    prisoner = dict((k.lower(), v) for k, v in json.loads(prisoner).items())
                    if "prisoner" in prisoner and prisoner["prisoner"].isdigit():
                        prisoner = int(prisoner["prisoner"])
                        if 1 <= prisoner <= len(PRISONERS):
                            return wrap_response(200,
                                                 {'text': PRISONERS[prisoner]})
                        return wrap_response(400, {
                            'text': "Prisoner was not found"
                        })
                    return wrap_response(451, {
                        'text': "Unavailable For Legal Reasons, Stop messing around!"
                    })
                return wrap_response(200, {
                    "text": "Hello Warden"
                })
            return wrap_response(401, {
                'text': 'Authentication Failed'
            })


def wrap_response(code, res=None):
    return {
        'statusCode': code,
        'body': json.dumps(res) if res is not None else None,
        'headers': {
            "Content-Type": "application/json",
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Methods": "OPTIONS, POST",
            "Access-Control-Allow-Headers": "Access-Control-Allow-Origin, "
                                            "Access-Control-Allow-Methods, Content-Type, authorization"
        }
    }


def lambda_handler(event, _):
    if event["httpMethod"] == "OPTIONS":
        return wrap_response(200)
    try:
        if event["httpMethod"] == "POST":
            headers = {k.lower(): v for k, v in event["headers"].items()}
            if "body" in event:
                return get_user(headers["authorization"], event["body"])
            else:
                return get_user(event["authorization"])
    except:
        return wrap_response(500, {
            'text': 'Not Authorized!'
        })
