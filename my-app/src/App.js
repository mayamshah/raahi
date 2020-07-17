import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import axios from "axios";
import Map from "./Map.js"
import MyMapComponent from "./MyMapComponent.js"
import StravaMapComponent from "./StravaMapComponent.js"
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';

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
  const [currentIndex, setCurrentIndex] = useState(0);
  const [error, setError] = useState("");
  const [open, setOpen] = useState(false);
  const [show, setShow] = useState(false);
  const [showStrava, setStravaShow] = useState(false);
  const [path, setPath] = useState([40.443659, -79.944641, 40.443659, -79.944641, 40.443659, -79.944641, 40.443659, -79.944641]);
  const [start, setStart] = useState([0,0]);
  const [end, setEnd] = useState([0,0]);



  function getResult() {
     const request = {
        address: currentAddress,
        distance: currentMiles
      };

      axios.post(endpoint, request)
      .then(res => {
        console.log(res);
        console.log(res.data);
        if (res.data.Error == '') {
          console.log("No error")
          setResult(res.data.Path)
          setResDist(res.data.Distance)
          setCurrentIndex(0)
          setShow(true)
        } else {
          setError(res.data.Error)
          setOpen(true)
        }
        
        console.log(res.data.Path)
      })


  }


  const stravaendpoint = "http://localhost:8080/api/executestrava"

  function getStravaResult() {
     const request = {
        address: currentAddress,
        distance: currentMiles
      };

      axios.post(stravaendpoint, request)
      .then(res => {
        console.log(res);
        console.log(res.data);
        if (res.data.Error == '') {
          setPath(res.data.Path)
          console.log(res.data.Path)
          setStart(res.data.Start)
          setEnd(res.data.End)
          setStravaShow(true)
        } else {
          setError(res.data.Error)
          setOpen(true)
        }
      })
  }

  function handleAddress(e) {
    setShow(false)
    setStravaShow(false)
    setCurrentAddress(e.target.value)
  }

  function handleMiles(e) {
    setShow(false)
    setStravaShow(false)
    setCurrentMiles(e.target.value)
  }

  function nextResult() {
    if (currentIndex < 2) {
      setCurrentIndex(currentIndex+1)
    }
  }

  function prevResult() {
    if (currentIndex > 0) {
      setCurrentIndex(currentIndex - 1)
    }
  }

  const handleClose = () => {
    setOpen(false);
  };


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
        <Button className={classes.button} onClick={() => nextResult()} variant="contained" color="primary" >
        NextResult
        </Button>
        <Button className={classes.button} onClick={() => prevResult()} variant="contained" color="primary" >
        PreviousResult
        </Button>
        <Button className={classes.button} onClick={() => getStravaResult()} variant="contained" color="primary" >
        Routes Near Me
        </Button>
        {show &&<Typography className={classes.heading} component="h1" variant="h6" gutterBottom>
          {resDist[currentIndex].toString() + " mile route below"}
        </Typography>
        }
        {show && <MyMapComponent org = {result[currentIndex]} />}
        {showStrava && <StravaMapComponent path = {path} start = {start} end = {end} />}
        <Dialog
          open={open}
          onClose={handleClose}
          aria-labelledby="alert-dialog-title"
          aria-describedby="alert-dialog-description"
        >
          <DialogTitle id="alert-dialog-title">{"An error occured"}</DialogTitle>
          <DialogContent>
            <DialogContentText id="alert-dialog-description">
              {error}
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose} color="primary">
              Close
            </Button>
          </DialogActions>
      </Dialog>

    </div>
    );
}

function App() {

  const classes = useStyles();
  return (
    <div className={classes.paper}> 
      <Display />
    </div>
  );
}


export default App;
