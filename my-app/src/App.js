import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import axios from "axios";
import MyMapComponent from "./MyMapComponent.js"
import DirectionsViewer from "./DirectionsViewer.js"
import NewDirectionsViewer from "./NewDirectionsViewer.js"
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import CssBaseline from '@material-ui/core/CssBaseline';
import Paper from '@material-ui/core/Paper';
import Fade from '@material-ui/core/Fade';
import DirectionsNew from "./NewMap.js";
import {LoadScript} from '@react-google-maps/api';
import Grid from '@material-ui/core/Grid';
import getDirections from './DirectionsExporter.js'
import './title.css';



let endpoint = "http://localhost:8080/api/execute";

const useStyles = makeStyles((theme) => ({

  overall_layout: {
    display: 'flex',
    flexWrap: "wrap",
    // [theme.breakpoints.down(1100)]: {
    //   display: "",
    //   flexWrap: "nowrap",
    // }
  },

  layout: {
    width: 400,
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

  directions_layout:{
    width: 600,
    marginLeft: 'auto',
    marginRight: 'auto',
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
  root_grid: {
    flexGrow: 1,
  },
  paper_grid: {
    padding: theme.spacing(1),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
}));

function NestedGrid(props) {
  const classes = useStyles();
  const results = props.response
  const mapClicked = props.onClick

  const Something = React.memo(props =>
  {
    console.log("render")
    console.log(props.response)
    return <DirectionsNew response = {props.response} onClick={props.onClick}/>;
  });

  return (
    <div className={classes.root_grid}>
      <Grid container spacing={1}>
      {props.response.map((res, index) => (
        <Grid item xs={4}>
          <Something response = {res} onClick= {() => mapClicked(index)}/>
        </Grid>
      ))}
      </Grid>
    </div>
  );
}

function Display() {
  const classes = useStyles();
  const [result, setResult] = useState([]);
  const [currentAddress, setCurrentAddress] = useState("");
  const [currentMiles, setCurrentMiles] = useState("0");
  const [currentIndex, setCurrentIndex] = useState(0);
  const [error, setError] = useState("");
  const [open, setOpen] = useState(false);
  const [show, setShow] = useState(false);
  const [showDirections, setShowDirections] = useState(false);
  const [gridView, setGridView] = useState(false);
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
          console.log(res.data.Results[currentIndex].Directions);
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
    setShowDirections(false)
    setCurrentAddress(e.target.value)
  }

  function handleMiles(e) {
    setShow(false)
    setShowDirections(false)
    setCurrentMiles(e.target.value)
  }

  function nextResult() {
    if (currentIndex < result.length - 1) {
      setCurrentIndex(currentIndex+1)
      setShowDirections(false)
    }
  }

  function prevResult() {
    if (currentIndex > 0) {
      setCurrentIndex(currentIndex - 1)
      setShowDirections(false)
    }
  }

  function buttonShowDirections() {
    console.log("show directions");
    setShowDirections(true)
  }

  function makeWaypoints(org) {
    var waypoints = []
    var i = 0
    for (i = 0; i < org.length && i < 46; i = i + 2) {
      waypoints.push({latitude: org[i], longitude: org[i+1]})
    }

    return waypoints
  }
  function buttonExportDirections() {
    const data = {
      source: {
        latitude: -33.8356372,
        longitude: 18.6947617
      },
      destination: {
        latitude: -33.8600024,
        longitude: 18.697459
      },
      params: [
        {
          key: "travelmode",
          value: "walking"        // may be "walking", "bicycling" or "transit" as well
        },
        {
          key: "dir_action",
          value: "navigate"       // this instantly initializes navigation using the given travel mode
        }
      ], 
      waypoints: []}
    data.source.latitude = result[currentIndex].Org[0]
    data.source.longitude = result[currentIndex].Org[1]
    data.destination.latitude = result[currentIndex].Dest[0]
    data.destination.longitude = result[currentIndex].Dest[1]
    data.waypoints = makeWaypoints(result[currentIndex].Path)
    getDirections(data)
  }

  function buttonGridView() {
    console.log("gridView");
    setGridView(true);
    setShow(false);
  }

  function normalView() {
    setShow(true);
    setGridView(false);
  }

  function mapClicked(index) {
    console.log("map clicked");
    console.log(index)
    setCurrentIndex(index)
    setGridView(false)
    setShow(true)
  }

  const handleClose = () => {
    setOpen(false);
  };


  const Something = React.memo(props =>
  {
    console.log("render")
    console.log(props.response)
    return <DirectionsNew response = {props.response} onClick= {props.onClick} />;
  });

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
    <LoadScript
        googleMapsApiKey="AIzaSyB32cCcL4gD_WIYPP6dAVSprY_QYE3arsk"
      >
    {show &&
      <main className={classes.map_layout}>
        <Fade in={show}>
            <Paper className={classes.paper}>
            <Something response = {result[currentIndex]} onClick= {() => mapClicked()}/>
            <Typography className={classes.heading} component="h1" variant="h6" gutterBottom>
            {result[currentIndex].Distance.toString() + " mile route"}
            </Typography>
            <Button className={classes.button} onClick={() => nextResult()} variant="contained" color="primary" >
            Next
            </Button>
            <Button className={classes.button} onClick={() => prevResult()} variant="contained" color="primary" >
            Previous
            </Button>
            <Button className={classes.button} onClick={() => buttonShowDirections()} variant="contained" color="primary" >
            Show Directions
            </Button>
            <Button className={classes.button} onClick={() => buttonExportDirections()} variant="contained" color="primary" >
            Export Directions
            </Button>
            { false && <Button className={classes.button} onClick={() => buttonGridView()} variant="contained" color="primary" >
            GridView
            </Button> }
            </Paper>
        </Fade>
      </main>
    }
    {false && 
      <main className={classes.map_layout}>
        <Paper className={classes.paper}>
        <NestedGrid response= {result} onClick={mapClicked}/>
        <Button className={classes.button} onClick={() => normalView()} variant="contained" color="primary" >
            Normal View
        </Button>
        </Paper>
      </main>
    }
    {showDirections &&
      <main className={classes.directions_layout}>
          <Paper className={classes.paper}>
            <NewDirectionsViewer steps={result[currentIndex].Directions}/>
          </Paper>
       </main>
     }
    </LoadScript>
    <main>
      {loading &&
        <Typography className={classes.heading} component="h1" variant="h6" gutterBottom>
          Loading
        </Typography>
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
