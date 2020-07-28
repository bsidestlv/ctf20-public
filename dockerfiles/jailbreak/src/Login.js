import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';
import FormControl from 'react-bootstrap/FormControl';
import React, { Component } from 'react';
import Image from 'react-bootstrap/Image'
import Alert from 'react-bootstrap/Alert';
import jwtDecode from "jwt-decode";
import { Redirect } from "react-router-dom";


class Login extends Component {

    constructor(props) {
        super(props);
        this.state = {
            username: '',
            password: '',
            isSent: false,
            error: null,
            IdToken: null,
            AccessToken: null
        };
        this.signInHandle = this.signInHandle.bind(this);
    }

    componentWillMount() {
        let AccessToken = localStorage.getItem('AccessToken');
        AccessToken = AccessToken ? AccessToken : 'unknown';
        let IdToken = localStorage.getItem('IdToken');
        IdToken = IdToken ? IdToken : 'unknown';
        if(AccessToken !== 'unknown' && IdToken !== 'unknown') {
            try {
                IdToken = jwtDecode(IdToken);
            } catch (e) {
                IdToken = 'unknown';
            }
        }
        this.setState({IdToken: IdToken});
        this.setState({AccessToken: AccessToken});
    }

    signIn(username, password) {
        let data = {username, password};

        fetch('https://jailbreak-auth.ctf.bsidestlv.com/Authenticate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        })
            .then(response => response.json())
            .then(data => {
                this.setState({'isSent': false});
                if(data.isSuccees.toString() === 'true') {
                    localStorage.setItem('AccessToken', data.AccessToken);
                    localStorage.setItem('IdToken', data.IdToken);
                    return window.location.href = '/';
                }
                return this.setState({'error': data.text});
            })
            .catch((data) => {
                this.setState({'isSent': false});
                if(!data.isSuccess) {
                    this.setState({'error': data.text})
                }
            });
    }

    signInHandle() {
        if(!this.state.username || !this.state.password) {
            return this.setState({'error': 'Please fill username and password!'})
        }
        this.signIn(this.state.username, this.state.password);
        this.setState({'isSent': true})
    }

    render() {
        if(this.state.IdToken !== 'unknown' && this.state.AccessToken) {
            return <Redirect to="/" />
        } else {
            return (
                <div>
                    <Image src={'/images/1.jpg'}
                           style={{'position': 'absolute', 'zIndex': '-1', 'width': '100%', 'opacity': 0.8}} />
                    <div style={{'text-align': 'center', 'zIndex': '9999'}}>
                        <Form style={{'width': '500px', 'marginLeft': 'auto', 'marginRight': 'auto'}}>
                            <div style={{'paddingTop':'50%'}}>
                                <InputGroup size="sm" className="mb-3">
                                    <InputGroup.Prepend>
                                        <InputGroup.Text id="inputGroup-sizing-sm">Username</InputGroup.Text>
                                    </InputGroup.Prepend>
                                    <FormControl aria-label="Username" aria-describedby="inputGroup-sizing-sm" onChange={(event) => this.setState({username: event.target.value})} />
                                </InputGroup>

                                <InputGroup size="sm" className="mb-3">
                                    <InputGroup.Prepend>
                                        <InputGroup.Text id="inputGroup-sizing-sm">Password</InputGroup.Text>
                                    </InputGroup.Prepend>
                                    <FormControl aria-label="Password" type={'password'} aria-describedby="inputGroup-sizing-sm" onChange={(event) => this.setState({password: event.target.value})} />
                                </InputGroup>

                                <Button variant="dark" type="button" onClick={() => {this.signInHandle()}} disabled={this.state.isSent}>
                                    {this.state.isSent ? "Loading.." : "Login"}
                                </Button>
                                {this.state.error && (
                                    <Alert key='1' variant='danger' style={{'marginLeft': 'auto', 'marginTop': '10px', 'marginRight': 'auto', 'textAlign': 'center'}}>
                                        <b>ERROR!</b><br />
                                        {this.state.error}
                                    </Alert>
                                )}
                            </div>
                        </Form>
                    </div>
                </div>
            );
        }
    }
}

export default Login;
