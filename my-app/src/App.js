import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import axios from "axios";
import MyMapComponent from "./MyMapComponent.js"
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import CssBaseline from '@material-ui/core/CssBaseline';
import Paper from '@material-ui/core/Paper';
import Fade from '@material-ui/core/Fade';
import './title.css';



let endpoint = "http://localhost:8080/api/execute";

const useStyles = makeStyles((theme) => ({

  overall_layout: {
    display: 'flex',
    flexWrap: "wrap",
  },

  layout: {
    width: 500,
    marginLeft: 'auto',
    marginRight: 'auto',
    // [theme.breakpoints.down(1010)]: {
    //   marginLeft: 'auto',
    //   marginRight: 'auto',
    // },
  },

  map_layout: {
    width: 600,
    marginLeft: 'auto',
    marginRight: 'auto',
    // [theme.breakpoints.down(1000)]: {
    //   marginLeft: 'auto',
    //   marginRight: 'auto',
    // },
  },

  paper: {
    marginTop: theme.spacing(3),
    marginLeft: 'auto',
    marginRight: 'auto',
    padding: theme.spacing(2),
  },

  paper_for_map: {
    height: `100%`,
  padding: theme.spacing(2),

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
  const [result, setResult] = useState([]);
  const [currentAddress, setCurrentAddress] = useState("");
  const [currentMiles, setCurrentMiles] = useState("0");
  const [currentIndex, setCurrentIndex] = useState(0);
  const [error, setError] = useState("");
  const [open, setOpen] = useState(false);
  const [show, setShow] = useState(false);
  const [loading, setLoading] = useState(false);


  function getResult() {
     const request = {
        address: currentAddress,
        distance: currentMiles
      };
      setLoading(true)

      axios.post(endpoint, request)
      .then(res => {
        console.log(res);
        console.log(res.data);
        setLoading(false)
        if (res.data.Error == '') {
          console.log("No error")
          setResult(res.data.Results)
          setCurrentIndex(0)
          setShow(true)
        } else {
          setError(res.data.Error)
          setOpen(true)
        }

        
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
          console.log("No error")
          setResult(res.data.Results)
          setCurrentIndex(0)
          setShow(true)
        } else {
          setError(res.data.Error)
          setOpen(true)
        }
      })
  }

  function handleAddress(e) {
    setShow(false)
    setCurrentAddress(e.target.value)
  }

  function handleMiles(e) {
    setShow(false)
    setCurrentMiles(e.target.value)
  }

  function nextResult() {
    if (currentIndex < result.length - 1) {
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
    <div className={classes.overall_layout}>
    <main className={classes.layout}>
      <Paper className={classes.paper}>
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
        <Button className={classes.button} onClick={() => getStravaResult()} variant="contained" color="primary" >
        Routes Near Me
        </Button>
      </Paper>
    </main>
    <main>
      {loading &&
        <Typography className={classes.heading} component="h1" variant="h6" gutterBottom>
          Loading
        </Typography>
      }
    </main>
    <main className={classes.map_layout}>
      {show &&
      <Fade in={show}>
          <Paper className={classes.paper}>
          <MyMapComponent response = {result[currentIndex]} />
          <Typography className={classes.heading} component="h1" variant="h6" gutterBottom>
          {result[currentIndex].Distance.toString() + " mile route"}
          </Typography>
          <Button className={classes.button} onClick={() => nextResult()} variant="contained" color="primary" >
          Next
          </Button>
          <Button className={classes.button} onClick={() => prevResult()} variant="contained" color="primary" >
          Previous
          </Button>
          </Paper>
      </Fade>
      }
    </main>
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

const overallStyles = makeStyles((theme) => ({
   
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

function App() {

  const classes = overallStyles();

  return (
    <React.Fragment>
      <CssBaseline />
        <Typography className="title.body" component="h1" variant="h4" align="center">
          Rahi
        </Typography>
        <Display/>
    </React.Fragment>
  );
}


export default App;
