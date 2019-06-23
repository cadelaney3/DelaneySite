import React, { Component, useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import './App.css';
import Header from './components/Header/Header';
import ChatHistory from './components/ChatHistory/ChatHistory';
import ChatInput from './components/ChatInput/ChatInput';
import Home from './components/Home/Home';
import SignIn from './components/SignIn/SignIn';
import { connect, sendMsg } from "./api";
import { createMuiTheme } from '@material-ui/core/styles';
import { MuiThemeProvider } from '@material-ui/core/styles';

const theme = createMuiTheme({
  palette: {
    type: 'dark',
  }
})
/*
class App extends Component {
  constructor(props) {
    super(props);
  }

  state = {
    loggedIn: false,
  }

  handleSignIn() {
    this.setState(prevState => ({ loggedIn: !prevState.loggedIn }));
  }

  render() {
    const [signedIn, setSignedIn] = useState(false);

    return (
      <Router>
        <MuiThemeProvider theme={theme} >
          <Header theme={theme} />
          <Switch>
            <Route expact path="/home" render={(props) => <HomePage {...props} loggedIn={signedIn} /> } />
            <Route path="/ws" render={(props) => <WS {...props} loggedIn={this.state.loggedIn} /> } />
            <Route 
              path="/signin" 
              render={(props) => <SignIn
                {...props}
                loggedIn={signedIn} //{this.state.loggedIn} 
                handleSignIn={setSignedIn} //{this.handleSignIn}  
                />            
              } 
            />
          </Switch>
        </MuiThemeProvider>
      </Router>
    );
  }
}
*/
export default function App() {
  const [loggedIn, setLoggedIn] = useState(false);

  useEffect(() => {
    console.log("login: ", loggedIn);
    sessionStorage.setItem('loggedIn', loggedIn);
  });

  const handleSignIn = () => {
    setLoggedIn(!loggedIn);
  }

  return (
    <Router>
      <MuiThemeProvider theme={theme} >
        <Header theme={theme} />
        <Switch>
          <Route expact path="/home" render={(props) => <HomePage {...props} loggedIn={loggedIn} /> } />
          <Route path="/ws" render={(props) => <WS {...props} loggedIn={loggedIn} /> } />
          <Route 
            path="/signin" 
            render={(props) => <SignIn
              {...props}
              loggedIn={loggedIn} //{this.state.loggedIn} 
              handleSignIn={handleSignIn} //{this.handleSignIn}  
              />            
            } 
          />
        </Switch>
      </MuiThemeProvider>
    </Router>
  );


}

class HomePage extends Component {
  constructor(props) {
    super(props);

    console.log(sessionStorage.getItem('loggedIn'));

    this.state = {
      error: null,
      isLoaded: false,
      items: []
    };
  }

  componentDidMount() {
    //fetch("http://localhost:8080/home")
    fetch("http://172.17.21.104:8080/home")
    .then(res => res.json())
    .then(
      (result) => {
        this.setState({
          isLoaded: true,
          items: result
        });
      },
      (error) => {
        this.setState({
          isLoaded: true,
          error
        });
      }
    )
  }

  render() {
    const { error, isLoaded, items } = this.state;
    if (error) {
      return <div>Error: {error.message}</div>;
    } else if (!isLoaded) {
      return <div>Loading...</div>;
    } else {
      return (
        <div className="HomePage">
          <Home items={items} />
        </div>
      );
    }
  }

}

class WS extends Component {
  constructor(props) {
    super(props);

    this.state = {
      prevState: null,
      chatHistory: []
    }
  }

  componentDidMount() {
    connect((msg) => {
      console.log("New Message")
      this.setState(prevState => ({
        chatHistory: [...this.state.chatHistory, msg]
      }))
      console.log(this.state);
    });
  }

  send(event) {
    if(event.keyCode === 13) {
      sendMsg(event.target.value);
      event.target.value = "";
    }
  }

  render() {
    return (
      <div className="WS">
        <ChatHistory chatHistory={this.state.chatHistory} />
        <ChatInput send={this.send} />
      </div>
    );
  }
}

// export default App;
