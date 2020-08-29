import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import axios from "axios";
import MyMapComponent from "./MyMapComponent.js"
import DirectionsViewer from "./DirectionsViewer.js"
import NewDirectionsViewer from "./NewDirectionsViewer.js"
import TrailViewer from "./TrailViewer.js"
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
import Toolbar from '@material-ui/core/Toolbar';
import AppBar from '@material-ui/core/AppBar';
import { Box } from '@material-ui/core';
import IconButton from '@material-ui/core/IconButton';
import InfoIcon from '@material-ui/icons/Info';
import MenuIcon from '@material-ui/icons/Menu';
import DirectionsIcon from '@material-ui/icons/Directions';
import DirectionsRunIcon from '@material-ui/icons/DirectionsRun';
import TimelineIcon from '@material-ui/icons/Timeline';
import LinkedInIcon from '@material-ui/icons/LinkedIn';
import GitHubIcon from '@material-ui/icons/GitHub';
import MailIcon from '@material-ui/icons/Mail';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import Tooltip from '@material-ui/core/Tooltip';
import Drawer from '@material-ui/core/Drawer';
import Collapse from '@material-ui/core/Collapse';
import ListItemText from '@material-ui/core/ListItemText';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import clsx from 'clsx';
import CircularProgress from '@material-ui/core/CircularProgress';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import CardMedia from '@material-ui/core/CardMedia';
import Container from '@material-ui/core/Container';
import Divider from '@material-ui/core/Divider';
import Slide from '@material-ui/core/Slide';
import ArrowForwardIcon from '@material-ui/icons/ArrowForward';
import { createMuiTheme, ThemeProvider } from '@material-ui/core/styles';
import Checkbox from '@material-ui/core/Checkbox';
import { green, orange } from '@material-ui/core/colors';








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
  container: {
    display: 'flex',
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

  map_button_layout: {
    width: 400
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
    height: 400,
    width: 200,
  },

  heading: {
    margin: theme.spacing(0,0,0,0),
    color: 'black',
    height: 48,
    align: 'center',
  },
  textField: {
    width: '25ch',
    margin: theme.spacing(0,0,4,0),
  },
  button: {
    margin: theme.spacing(0,1,1,0),
  },
  np_button: {
    margin: theme.spacing(2,1,1,0)

  },
  root_grid: {
    flexGrow: 1,
  },
  paper_grid: {
    padding: theme.spacing(1),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  largeIcon: {
    width: 60,
    height: 60,
  },

  // button_box: {
  //   align: 'center',
  // },
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
  const [trails, setTrails] = useState([]);
  const [origin, setOrigin] = useState([]);
  const [currentAddress, setCurrentAddress] = useState("");
  const [currentMiles, setCurrentMiles] = useState("0");
  const [currentIndex, setCurrentIndex] = useState(0);
  const [error, setError] = useState("");
  const [open, setOpen] = useState(false);
  const [show, setShow] = useState(false);
  const [showTrail, setShowTrail] = useState(false);
  const [showDirections, setShowDirections] = useState(false);
  const [gridView, setGridView] = useState(false);
  const [loading, setLoading] = useState(false);
  const [menu, setMenu] = useState(true);



  function getResult() {
     const request = {
        address: currentAddress,
        distance: currentMiles
      };
      setLoading(true)
      setShowTrail(false)

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
          setShowTrail(false)
        } else {
          setError(res.data.Error)
          setOpen(true)
        }

        
      })


  }

  const trailendpoint = "http://localhost:8080/api/executetrail"

  function getTrailResult() {
     const request = {
        address: currentAddress,
        distance: currentMiles
      };
      setLoading(true)
      setShow(false)
      setShowDirections(false)

      axios.post(trailendpoint, request)
      .then(res => {
        console.log(res);
        console.log(res.data);
        setLoading(false)
        if (res.data.Error == '') {
          console.log("No error")
          setTrails(res.data.Results)
          setOrigin(res.data.Origin)
          setShowTrail(true)
          setShow(false)
          // setCurrentIndex(0)
          // setShow(true)
        } else {
          setError(res.data.Error)
          setOpen(true)
        }
      })
  }

  function handleAddress(e) {
    setShow(false)
    setShowTrail(false)
    setShowDirections(false)
    setCurrentAddress(e.target.value)
  }

  function handleMiles(e) {
    setShow(false)
    setShowTrail(false)
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
        <Typography className={classes.heading} component="h1" variant="h6">
          What is your location?
        </Typography>
        <form className={classes.textField} noValidate>
          <TextField id="standard-basic" label="Address" onChange={e => handleAddress(e)} onKeyPress={(ev) => {
                                                                                            console.log(`Pressed keyCode ${ev.key}`);
                                                                                            if (ev.key === 'Enter') {
                                                                                              // Do code here
                                                                                              ev.preventDefault();
                                                                                            }
                                                                                          }}/>
        </form>
        <Typography className={classes.heading} component="h1" variant="h6">
          What is your desired distance?
        </Typography>
        <form className={classes.textField} noValidate>
          <TextField id="standard-basic" label="Miles" onChange={e => handleMiles(e)}/>
        </form>
        <Box alignItems="center" justifyContent="center" display="flex" >
          <Button className={classes.button} onClick={() => getResult()} variant="contained" color="primary" fullWidth>
          Generate Routes
          </Button>
          <Button className={classes.button} onClick={() => getTrailResult()} variant="contained" color="primary" fullWidth>
          Trails Near Me
          </Button>
        </Box>
        
      </Paper>
       <Box alignItems="center" justifyContent="center" display="flex" paddingTop={5}>
        {loading &&
          <CircularProgress />
        }
      </Box>
    </main>
    <LoadScript
        googleMapsApiKey="AIzaSyB32cCcL4gD_WIYPP6dAVSprY_QYE3arsk"
      >
    {show &&
      <main className={classes.map_layout}>
        <Fade in={show}>
            <Paper className={classes.paper}>
            <Box alignItems="center" justifyContent="center" display="flex" >
                  <Paper className={classes.paper_for_map}> 
                  <Grid container direction="column" spacing={4}>
                      <Grid item xs={12} zeroMinWidth>
                      <Typography component="h1" variant="overline" align="center" >
                        {"Showing result: " + (currentIndex + 1) + "/" + result.length}
                      </Typography>
                      </Grid>

                      <Grid item xs={12} zeroMinWidth>
                      <Typography component="h1" variant="h6" align="center">
                          {Number(result[currentIndex].Distance).toFixed(2).toString() + " mile route"}
                      </Typography>
                      </Grid>

                      <Grid item xs={12}>
                        <Grid container alignItems="flex-start" justify="center" direction="row">
                        <Tooltip title="Directions">
                              <IconButton onClick={() => buttonShowDirections()} >
                                <DirectionsRunIcon fontSize="large"/>
                              </IconButton>
                        </Tooltip>
                        </Grid>
                       </Grid> 

                      <Grid item xs={12}>
                        <Grid container alignItems="flex-start" justify="center" direction="row">
                        <Tooltip title="Export Route">
                              <IconButton onClick={() => buttonExportDirections()}>
                                <DirectionsIcon fontSize="large"/>
                              </IconButton>
                        </Tooltip>
                        </Grid>
                      </Grid>

                      <Grid item xs={12}>
                        <Grid container alignItems="flex-start" justify="center" direction="row">
                        <Tooltip title="Trails Near Me">
                              <IconButton onClick={() => getTrailResult()}>
                                <TimelineIcon fontSize="large"/>
                              </IconButton>
                        </Tooltip>
                        </Grid>
                      </Grid>

                      


                  </Grid>
                  </Paper>
            <main className={classes.map_button_layout}>
              <Something response = {result[currentIndex]} onClick= {() => mapClicked()}/>
              <main className={classes.button_box}>
                <Box alignItems="center" justifyContent="center" display="flex" height={50}>
                  <Button className={classes.np_button} onClick={() => prevResult()} variant="contained" color="primary" fullWidth>
                  Previous
                  </Button>
                  <Button className={classes.np_button} onClick={() => nextResult()} variant="contained" color="primary" fullWidth>
                  Next
                  </Button>
                </Box>
              </main>
            </main>
            </Box>
            
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
    {showTrail && 
      <Paper className={classes.paper}>
        <TrailViewer trails={trails} origin={origin}/>
      </Paper>
    }
    {showDirections &&
      <main className={classes.directions_layout}>
          <Paper className={classes.paper}>
            <NewDirectionsViewer steps={result[currentIndex].Directions}/>
          </Paper>
       </main>
     }
    </LoadScript> 
    
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
    color: 'white',
    height: 48,
    align: 'left',
  },
  subtext: {
    margin: theme.spacing(0,0,0,3),
    color: 'white',
    height: 20,
    align: 'left',
  },
  textField: {
    width: '25ch',
    margin: theme.spacing(0,0,3,3),
  },
  button: {
    margin: theme.spacing(3,0,0,3),
  },
  root: {
    display: 'flex',
  },
  appBar: {
    transition: theme.transitions.create(['margin', 'width'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
  },
  appBarShift: {
    width: `calc(100% - ${drawerWidth}px)`,
    marginLeft: drawerWidth,
    transition: theme.transitions.create(['margin', 'width'], {
      easing: theme.transitions.easing.easeOut,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  menuButton: {
    marginRight: theme.spacing(2),
  },
  hide: {
    display: 'none',
  },
  drawer: {
    width: drawerWidth,
    flexShrink: 0,
  },
  drawerPaper: {
    width: drawerWidth,
  },
  drawerHeader: {
    display: 'flex',
    alignItems: 'center',
    padding: theme.spacing(0, 1),
    // necessary for content to be below app bar
    ...theme.mixins.toolbar,
    justifyContent: 'flex-end',
  },
  content: {
    flexGrow: 1,
    padding: theme.spacing(3),
    transition: theme.transitions.create('margin', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    marginLeft: -drawerWidth,
  },
  contentShift: {
    transition: theme.transitions.create('margin', {
      easing: theme.transitions.easing.easeOut,
      duration: theme.transitions.duration.enteringScreen,
    }),
    marginLeft: 0,
  },

  paper: {
      marginTop: theme.spacing(3),
      marginLeft: 'auto',
      marginRight: 'auto',
      padding: theme.spacing(2),
      width: 600,
    },
}));

const drawerWidth = 240;

const albumStyles = makeStyles((theme) => ({

    heroContent: {
       [theme.breakpoints.up(1000)]: {
        margin: theme.spacing(5,10,0,10),
        },
        margin: theme.spacing(5,5,0,5),
      },
    textBox: {
      width: 500
    },

    root: {
    display: 'flex',
    flexDirection: 'column',
      },
  media: {
    height: 300,
  },
  controls: {
    display: 'flex',
    alignItems: 'center',
    paddingLeft: theme.spacing(1),
    paddingBottom: theme.spacing(1),
  },

  cardGrid: {
    paddingLeft: theme.spacing(5),
    paddingRight: theme.spacing(5),
    paddingBottom: theme.spacing(10),
  },

    overall_layout: {
    display: 'flex',
    flexWrap: "wrap",
    // [theme.breakpoints.down(1100)]: {
    //   display: "",
    //   flexWrap: "nowrap",
    // }
  },
  container: {
    display: 'flex',
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

}));

function About() {
  const classes = albumStyles()
  return (
    <div className={classes.heroContent}>

        <Typography component="h1" variant="h3" align="center" color="textPrimary" gutterBottom>
          Our Name
        </Typography>
        <Box p={5}>
        <Typography variant="h5" align="center" color="textSecondary" paragraph gutterBottom>
          Raahi means traveler in Hindi. We are all travelers on the path of life, and hopefully this app can find you some more paths on which to travel.
        </Typography>
        </Box>
        <Divider />
        <Box p={5}>
         <Typography component="h1" variant="h3" align="center" color="textPrimary" gutterBottom>
            Our Team
        </Typography>
        </Box>
       
    <Container className={classes.cardGrid} maxWidth="md">
    <Grid container spacing={4}>
      <Grid item xs={6} >
      <Card className={classes.root}>
          <CardMedia
            className={classes.media}
            image={require('./images/maya.jpg')}
          />
          <CardContent>
            <Typography gutterBottom variant="h5" component="h2">
              Maya Shah
            </Typography>
            <Typography variant="body2" color="textSecondary" component="p">
              Maya is from Saratoga, California majoring in Business Administration at Carnegie Mellon University pursuing an additional major in Computer Science. 
            </Typography>
          </CardContent>
        <div className={classes.controls}>
          <IconButton aria-label="linkedin" href="https://www.linkedin.com/in/mayamshah/">
            <LinkedInIcon/>
          </IconButton>
          <IconButton aria-label="github" href="https://github.com/mayamshah">
            <GitHubIcon/>
          </IconButton>
          <IconButton aria-label="mail" href="mailto:mayashah@andrew.cmu.edu">
            <MailIcon/>
          </IconButton>
        </div>
      </Card>
      </Grid>
      <Grid item xs={6} >
      <Card className={classes.root}>
          <CardMedia
            className={classes.media}
            image={require('./images/agam.jpg')}
          />
          <CardContent>
            <Typography gutterBottom variant="h5" component="h2">
              Agam Kumar
            </Typography>
            <Typography variant="body2" color="textSecondary" component="p">
              Agam is from Edison, New Jersey majoring in Statistics and Machine Learning at Carnegie Mellon University pursuing an additional major in Computer Science.             </Typography>
          </CardContent>
        <div className={classes.controls}>
          <IconButton aria-label="linkedin" href="https://www.linkedin.com/in/agam-kumar/">
            <LinkedInIcon/>
          </IconButton>
          <IconButton aria-label="github" href="https://github.com/agam-kumar">
            <GitHubIcon/>
          </IconButton>
          <IconButton aria-label="mail" href="mailto:agamk@andrew.cmu.edu">
            <MailIcon/>
          </IconButton>
        </div>
      </Card>
    </Grid>
    </Grid>
    </Container>

     <Typography variant="h5" align="center" color="textSecondary" paragraph>
         The two developers met in college. They realized they worked well together and had a shared passion for software development. Raahi thus started out as just the idea to continue working together over the summer. In efforts of staying active over quarantine, the developers searched for routes of a given distance that they could run in their neighborhoods to no avail. They now present this app to aide those with similar interests.
        </Typography>

    </div>
    )
}

function App() {

  const classes = overallStyles();
    const [drawerState, setDrawerState] = useState(false);
    const [display, setDisplay] = useState(false);
    const [about, setAbout] = useState(false);
    const [intro, setIntro] = useState(true);
    const outerTheme = createMuiTheme({
  palette: {
    
  },
});

    function homeClicked() {
      setDisplay(true) 
      setAbout(false)
      setIntro(false)
      setDrawerState(false)
    }

    function aboutClicked() {
      setDisplay(false) 
      setAbout(true)
      setIntro(false)
      setDrawerState(false)
    }

    function instructionsClicked() {
      setDisplay(false) 
      setAbout(false)
      setIntro(true)
      setDrawerState(false)
    }

  return (
    <React.Fragment>
      <CssBaseline />
      <ThemeProvider theme={outerTheme}>
      <AppBar position="relative"
      className={clsx(classes.appBar, {
          [classes.appBarShift]: drawerState,
        })}
      >
        <Toolbar>
        {true && 
            <div>
            <IconButton align='right' aria-label="route info" onClick={ () => setDrawerState(true)}>
               <MenuIcon style={{ color: '#FFFFFF' }}/>
            </IconButton>
            <Drawer anchor={'left'} open={drawerState}className={classes.drawer} classes={{
          paper: classes.drawerPaper,
        }} onClose={() => setDrawerState(false)}>
            <List>
              <ListItem button onClick={() => homeClicked()}>
              <ListItemText primary="Home" />
              </ListItem>
              <ListItem button onClick={() => aboutClicked()}>
              <ListItemText primary="About" />
              </ListItem>
              <ListItem button onClick={() => instructionsClicked()}>
              <ListItemText primary="Instructions" />
              </ListItem>
            </List>
            </Drawer>
            </div>
          }
          <Typography className={classes.title} component="h1" variant="h5">
            Raahi
          </Typography>
        </Toolbar>
      </AppBar>
      <main
        
      >
       {intro && 
          <Fade in={true} timeout={{ enter: 3000 }}>
          <Paper elevation={4} className={classes.paper}>
            <Typography variant="h3" align="center">
            How to Use Raahi
            </Typography>
            <Box p={4}>
            <Divider/>
            </Box>
            <Typography component="h1" variant="h6" color="textSecondary" align="center">
            Enter a location and distance. Raahi will then generate routes near you to walk, run, or bike. The routes begin and end at your entered location and are at least your desired distance. Raahi also offers the option to find trails near you and directions to get to them.


            </Typography>
            <Typography component="h1" variant="h6" color="textSecondary">
            </Typography>
            <Typography component="h1" variant="h6" color="textSecondary">
            </Typography>

            <Slide direction="right" in={true} timeout={{ enter: 1000 }} mountOnEnter unmountOnExit>
            <Grid container alignItems="flex-start" justify="flex-end" direction="row">

              <Tooltip title="Go">
               <IconButton aria-label="mail" onClick={() => homeClicked()}>
                    <ArrowForwardIcon/>
              </IconButton>
              </Tooltip>
              </Grid>
            </Slide>
          </Paper>
        </Fade>
       }
       {display && <Display/>}
       {about && <About/>}
      </main>
      </ThemeProvider>
    </React.Fragment>
  );
}

export default App;
