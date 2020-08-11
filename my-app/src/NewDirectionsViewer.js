import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { GoogleMap, StreetViewPanorama} from '@react-google-maps/api'
import ListItem from '@material-ui/core/ListItem';
import Button from '@material-ui/core/Button';

function LatLng(lat, lng) {
  return {lat: lat, lng: lng};
}

function getHeading(lat1, lng1, lat2, lng2) {
  var point1 = LatLng(lat1, lng1)
  var point2 = LatLng(lat2, lng2)
  var heading = computeHeading(point1, point2)
  console.log(heading)
  return heading;
}

function StreetView(props) {
	const mapContainerStyle = {
	  height: "400px",
	  width: "100%"
	};

	const center = {
	  lat:  props.location[0],
	  lng: props.location[1]
  };
  
  // const headingInfo = {
  //   heading: getHeading(props.location[0], props.location[1], props.locEnd[0], props.locEnd[1]),
  //   pitch: 0
  // };

	return (
		<GoogleMap
			    id="circle-example"
			    mapContainerStyle={mapContainerStyle}
			    zoom={7}
			    center={center}
			  >
			    <StreetViewPanorama
			      position={center}
            visible={true}
            //pov={headingInfo}
			    />
	  	</GoogleMap>
	  	);
}
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


function NewDirectionsViewer(steps) {

	const classes = useStyles();
    const [currentIndex, setCurrentIndex] = useState(0);
    console.log(steps);

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
        <StreetView location={steps.steps[currentIndex].Loc} locEnd={steps.steps[currentIndex].EndLoc} />
  	</div>
  );

}

export default NewDirectionsViewer