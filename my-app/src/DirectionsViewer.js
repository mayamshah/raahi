import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import ListItem from '@material-ui/core/ListItem';
import Button from '@material-ui/core/Button';

const useStyles = makeStyles((theme) => ({

  heading: {
    margin: theme.spacing(3,0,0,3),
    color: 'black',
    height: 48,
  },
  button: {
    margin: theme.spacing(3,0,0,3),
  },
}));


const { compose, withProps } = require("recompose");
const {
  withScriptjs,
  withGoogleMap,
  GoogleMap,
  StreetViewPanorama,
  OverlayView,
} = require("react-google-maps");

const getPixelPositionOffset = (width, height) => ({
  x: -(width / 2),
  y: -(height / 2),
})

class StreetView extends React.Component {
  constructor(props){
    super(props)
  }

  
render() {
    const StreetViewPanorma = compose(
      withProps({
        googleMapURL: "https://maps.googleapis.com/maps/api/js?key=AIzaSyB32cCcL4gD_WIYPP6dAVSprY_QYE3arsk&v=3.exp&libraries=geometry,drawing,places",
        loadingElement: <div style={{ height: `100%` }} />,
        containerElement: <div style={{ height: `400px` }} />,
        mapElement: <div style={{ height: `100%` }} />,
        center: { lat: this.props.location[0], lng: this.props.location[1] },
      }),
      withScriptjs,
      withGoogleMap
    )(props =>
      <GoogleMap defaultZoom={8} defaultCenter={props.center}>
        <StreetViewPanorama defaultPosition={props.center} visible>
        </StreetViewPanorama>
      </GoogleMap>
    );
return (
        <div>
        <StreetViewPanorma />
        </div>
    )
  }
}


function DirectionsViewer(steps) {

    const classes = useStyles();
    const [currentIndex, setCurrentIndex] = useState(0);
    console.log(steps)

    function getNextStep() {
      if (currentIndex < steps.steps.length - 1) {
        setCurrentIndex(currentIndex + 1)
      }
    }

    function getPreviousStep() {
      if (currentIndex > 0) {
        setCurrentIndex(currentIndex - 1)
      }
    }

    return (
      <div>
        <ListItem>
        <div dangerouslySetInnerHTML={{__html: steps.steps[currentIndex].Instructions}} />
        </ListItem>
        <Button className={classes.button} onClick={() => getNextStep()} variant="contained" color="primary" >
          Next Step
        </Button>
        <Button className={classes.button} onClick={() => getPreviousStep()} variant="contained" color="primary" >
          Previous Step
        </Button>
        <StreetView location={steps.steps[currentIndex].Loc}/>
      </div>
      );
 
}

export default DirectionsViewer