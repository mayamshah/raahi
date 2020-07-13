import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import axios from "axios";
import Map from "./Map.js"
import MyMapComponent from "./MyMapComponent.js"
import StravaMapComponent from "./StravaMapComponent.js"

let endpoint = "http://localhost:8080/api/execute";

const useStyles = makeStyles((theme) => ({
  paper: {
    alignItems: 'center',
  },
  heading: {
    margin: theme.spacing(3,0,0,3),
    color: 'black',
    height: 48,
  },
  textField: {
    width: '25ch',
    margin: theme.spacing(0,0,3,3),
  },
  button: {
    margin: theme.spacing(3,0,0,3),
  },
}));

function Display() {
  const classes = useStyles();
  const [result, setResult] = useState([[40.443659, -79.944641, 40.443659, -79.944641, 40.443659, -79.944641, 40.443659, -79.944641], [40.443659, -79.944641, 40.443659, -79.944641, 40.443659, -79.944641, 40.443659, -79.944641], [40.443659, -79.944641, 40.443659, -79.944641, 40.443659, -79.944641, 40.443659, -79.944641],]);
  const [resDist, setResDist] = useState([0.0, 0.0, 0.0])
  const [currentAddress, setCurrentAddress] = useState("");
  const [currentMiles, setCurrentMiles] = useState("0");


  function getResult() {
     const request = {
        address: currentAddress,
        distance: currentMiles
      };

      axios.post(endpoint, request)
      .then(res => {
        console.log(res);
        console.log(res.data);
        setResult(res.data.Path)
        setResDist(res.data.Distance)
        console.log(res.data.Path)
      })
  }

  function handleAddress(e) {
    setCurrentAddress(e.target.value)
  }

  function handleMiles(e) {
    setCurrentMiles(e.target.value)
  }

  return (
    <div>
        <Typography className={classes.heading} component="h1" variant="h6" gutterBottom>
          What is your starting point?
        </Typography>
        <form className={classes.textField} noValidate>
          <TextField id="standard-basic" label="Address" onChange={e => handleAddress(e)}/>
        </form>
        <Typography className={classes.heading} component="h1" variant="h6" gutterBottom>
          How many miles?
        </Typography>
        <form className={classes.textField} noValidate>
          <TextField id="standard-basic" label="Miles" onChange={e => handleMiles(e)}/>
        </form>
        <Button className={classes.button} onClick={() => getResult()} variant="contained" color="primary" >
        Enter
        </Button>   
        <Typography className={classes.heading} component="h1" variant="h6" gutterBottom>
          {resDist[0].toString() + " mile route below"}
        </Typography>
        <MyMapComponent org = {result[0]} />
        <Typography className={classes.heading} component="h1" variant="h6" gutterBottom>
          {resDist[1].toString() + " mile route below"}
        </Typography>
        <MyMapComponent org = {result[1]} />
        <Typography className={classes.heading} component="h1" variant="h6" gutterBottom>
          {resDist[2].toString() + " mile route below"}
        </Typography>
        <MyMapComponent org = {result[2]} />
               
    </div>
    );
}


function StravaDisplay() {
  const classes = useStyles();
  const [path, setPath] = useState([40.443659, -79.944641, 40.443659, -79.944641, 40.443659, -79.944641, 40.443659, -79.944641]);
  const [start, setStart] = useState([0,0]);
  const [end, setEnd] = useState([0,0]);
  const [currentAddress, setCurrentAddress] = useState("");
  const [currentMiles, setCurrentMiles] = useState("0");

  const stravaendpoint = "http://localhost:8080/api/executestrava"

  function getResult() {
     const request = {
        address: currentAddress,
        distance: currentMiles
      };

      axios.post(stravaendpoint, request)
      .then(res => {
        console.log(res);
        console.log(res.data);
        setPath(res.data.Path)
        setStart(res.data.Start)
        setEnd(res.data.End)
        console.log(res.data.Path)
      })
  }

  function handleAddress(e) {
    setCurrentAddress(e.target.value)
  }

  function handleMiles(e) {
    setCurrentMiles(e.target.value)
  }

  return (
    <div>
        <Typography className={classes.heading} component="h1" variant="h6" gutterBottom>
          What is your starting point?
        </Typography>
        <form className={classes.textField} noValidate>
          <TextField id="standard-basic" label="Address" onChange={e => handleAddress(e)}/>
        </form>
        <Typography className={classes.heading} component="h1" variant="h6" gutterBottom>
          How many miles?
        </Typography>
        <form className={classes.textField} noValidate>
          <TextField id="standard-basic" label="Miles" onChange={e => handleMiles(e)}/>
        </form>
        <Button className={classes.button} onClick={() => getResult()} variant="contained" color="primary" >
        Enter
        </Button>   
        <StravaMapComponent path = {path} start = {start} end = {end} />
    </div>
    );
}

function App() {

  const classes = useStyles();
  return (
    <div className={classes.paper}> 
      <Display />
      <StravaDisplay />
    </div>
  );
}


export default App;
