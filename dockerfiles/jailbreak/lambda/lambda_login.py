import os
import hmac
import json
import boto3
import base64
import hashlib

# Global Variables
APP_CLIENT_ID = os.environ['APP_CLIENT_ID']
APP_CLIENT_SECRET = os.environ['APP_CLIENT_SECRET']


def get_secret_hash(username):
    digest = hmac.new(
        str(APP_CLIENT_SECRET).encode('utf-8'),
        msg=str(username + APP_CLIENT_ID).encode('utf-8'),
        digestmod=hashlib.sha256
    ).digest()
    return base64.b64encode(digest).decode()


def wrap_response(status_code, res=None):
    return {
        'statusCode': status_code,
        'body': json.dumps(res) if res is not None else None,
        'headers': {
            "Content-Type": "application/json",
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Methods": "OPTIONS, POST",
            "Access-Control-Allow-Headers": "Access-Control-Allow-Origin, Access-Control-Allow-Methods, Content-Type"
        }
    }


def authenticate_user(username, password):
    credentials = {
        'USERNAME': username.lower(),
        'SECRET_HASH': get_secret_hash(username.lower()),
        'PASSWORD': password
    }
    try:
        client = boto3.client('cognito-idp', region_name=os.environ['AWS_REGION'])
        response = client.initiate_auth(AuthFlow='USER_PASSWORD_AUTH',
                                        AuthParameters=credentials, ClientId=APP_CLIENT_ID)
    except Exception as e:
        if "attempts exceeded" in str(e).split("operation: ")[1]:
            return wrap_response(299, {
                'text': 'Violence is not the solution!',
                'isSuccees': False
            })
        return wrap_response(401, {
            'text': str(e).split("operation: ")[1],
            'isSuccees': False
        })
    if response and 'AccessToken' in response['AuthenticationResult'] and 'IdToken' in response['AuthenticationResult']:
        return wrap_response(200, {
            'isSuccees': True,
            'AccessToken': response['AuthenticationResult']['AccessToken'],
            'IdToken': response['AuthenticationResult']['IdToken']
        })
    return wrap_response(401, {
        'text': 'Authentication Failed',
        'isSuccees': False
    })


def lambda_handler(event, _):
    if event["httpMethod"] == "OPTIONS":
        return wrap_response(200)
    try:
        if event["httpMethod"] == "POST":
            auth_data = dict((k.lower(), v) for k, v in json.loads(event["body"]).items())
            if "username" not in auth_data or "password" not in auth_data:
                return wrap_response(500, {
                    'text': 'Please fill username and password.',
                    'isSuccees': False
                })
            return authenticate_user(auth_data["username"], auth_data["password"])
    except:
        return wrap_response(500, {
            'text': 'Please fill username and password.',
            'isSuccees': False
        })
    return wrap_response(2)
