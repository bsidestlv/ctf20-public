import Table from 'react-bootstrap/Table';
import Image from 'react-bootstrap/Image'
import Button from 'react-bootstrap/Button';
import Alert from 'react-bootstrap/Alert';
import { Redirect } from "react-router-dom";
import jwtDecode from 'jwt-decode'


import React, { Component } from 'react';

class PrisonManagement extends Component {
    constructor(props) {
        super(props);
        this.state = {
            AccessToken: null,
            IdToken: null,
            results: null,
            isSent: false,
        };
    }

    componentDidMount() {
        let AccessToken = localStorage.getItem('AccessToken');
        AccessToken = AccessToken ? AccessToken : 'unknown';
        let IdToken = localStorage.getItem('IdToken');
        IdToken = IdToken ? IdToken : 'unknown';
        if(AccessToken !== 'unknown' && IdToken !== 'unknown') {
            try {
                IdToken = jwtDecode(IdToken);
            } catch (e) {
                IdToken = 'unknown';
                AccessToken = 'unknown';
            }
        }
        this.setState({IdToken: IdToken});
        this.setState({AccessToken: AccessToken})
    }

    static logout() {
        localStorage.clear();
        window.location.href='/'
    }

    releasePrisoner(prisonerId) {
        if(!this.state.AccessToken || this.state === "unknown") {
            this.setState({results: "Try Harder!"})
        }
        this.setState({isSent: true});

        fetch('https://jailbreak-auth.ctf.bsidestlv.com/Validate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': this.state.AccessToken
            },
            body: JSON.stringify({'prisoner': prisonerId.toString()}),
        }).then(response => response.json()).then(res => {
            this.setState({results: res.text, isSent: false});

        })
            .catch(res => {
                this.setState({results: res.text, isSent: false})
            })

    }

    render() {
        if(!this.state.IdToken) {
            return (<div><h1>Loading...</h1></div>)

        }
        if(this.state.IdToken === 'unknown') {
            return <Redirect to="/login" />
        }
        if(this.state.IdToken['custom:isWarden'] === '1') {
            return (
                <div>
                    <Image src={'/images/2.jpg'}
                           style={{'position': 'absolute', 'zIndex': '-1', 'width': '100%', 'opacity': 0.8}}/>
                    <div style={{
                        'zIndex': `9999`,
                        'paddingTop': '250px',
                        'width': '50%',
                        'marginLeft': 'auto',
                        'marginRight': 'auto'
                    }}>
                        <Table striped bordered hover className={'table-dark'}>
                            <thead>
                            <tr>
                                <th>#</th>
                                <th>First Name</th>
                                <th>Last Name</th>
                                <th>Nickname</th>
                                <th>Management</th>
                            </tr>
                            </thead>
                            <tbody>
                            <tr>
                                <td>1</td>
                                <td>Artur</td>
                                <td>Isakhanyan</td>
                                <td>3viLMonk3y</td>
                                <td>
                                    <Button variant="outline-danger" onClick={() => {this.releasePrisoner(1)}} disabled={this.state.isSent}>
                                        {this.state.isSent ? "Loading.." : "Release Artur"}
                                    </Button>
                                </td>
                            </tr>
                            <tr>
                                <td>2</td>
                                <td>Nimrod</td>
                                <td>Levy</td>
                                <td>El3ct71k</td>
                                <td>
                                    <Button variant="outline-danger" onClick={() => {this.releasePrisoner(2)}} disabled={this.state.isSent}>
                                        {this.state.isSent ? "Loading.." : "Release Nimrod"}
                                    </Button>
                                </td>
                            </tr>
                            <tr>
                                <td>3</td>
                                <td>Groucho</td>
                                <td>Marx</td>
                                <td>J0hnSm1th</td>
                                <td>
                                    <Button variant="outline-danger" onClick={() => {this.releasePrisoner(3)}} disabled={this.state.isSent}>
                                        {this.state.isSent ? "Loading.." : "Release Groucho"}
                                    </Button>

                                </td>
                            </tr>
                            </tbody>
                        </Table>
                        {this.state.results && (
                            <Alert key='1' variant='success' style={{'marginLeft': 'auto', 'marginTop': '10px', 'marginRight': 'auto', 'textAlign': 'center'}}>
                                <b>Results</b><br />
                                {this.state.results}
                            </Alert>
                        )}
                        <div style={{"text-align": "center"}}>
                            <Button variant="danger" onClick={() => {PrisonManagement.logout()}} disabled={this.state.isSent}>
                                Logout
                            </Button>
                        </div>
                    </div>
                </div>
            )
        } else {
            return (
                <div>
                    <Image src={'/images/1.jpg'}
                           style={{'position': 'absolute', 'zIndex': '-1', 'width': '100%', 'opacity': 0.8}} />
                   <div style={{'zIndex': '9999', 'paddingTop': '20%'}}>
                       <Alert key='1' variant='danger' style={{'width': '30%', 'marginLeft': 'auto', 'marginRight': 'auto', 'textAlign': 'center'}}>
                           <b>FORBIDDEN!</b><br />
                           Only the warden is allowed!
                       </Alert>
                   </div>
                    <div style={{"text-align": "center"}}>
                        <Button variant="danger" onClick={() => {PrisonManagement.logout()}} disabled={this.state.isSent}>
                            Logout
                        </Button>
                    </div>
                </div>
            )
        }
    }
}


export default PrisonManagement;
