/*global google*/
import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { GoogleMap, StreetViewPanorama} from '@react-google-maps/api'
import ListItem from '@material-ui/core/ListItem';
import Button from '@material-ui/core/Button';
import { Box } from '@material-ui/core';
import Paper from '@material-ui/core/Paper';
import List from '@material-ui/core/List';





function getHeading(lat1, lng1, lat2, lng2) {
  var point1 = new google.maps.LatLng(lat1, lng1)
  var point2 = new google.maps.LatLng(lat2, lng2)
  var heading = google.maps.geometry.spherical.computeHeading(point1, point2)
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
  
  const headingInfo = {
    heading: getHeading(props.location[0], props.location[1], props.locEnd[0], props.locEnd[1]),
    pitch: 0
  };

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
            pov={headingInfo}
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
  np_button: {
    margin: theme.spacing(0,1,1,0)
  },
   paper: {
    height: 400,
    width: 200,
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

    function changeStep(index){
      setCurrentIndex(index)
    }

return (

	<div>
        <Box alignItems="center" justifyContent="center" display="flex" >
          {steps.steps && 
            <Paper style={{maxHeight: 400, overflow: 'auto', height:400}}>
            <List >
             {steps.steps.map((step, index) => (
              <ListItem button selected={index === currentIndex} onClick={() => changeStep(index)}>
               <div dangerouslySetInnerHTML={{__html: step.Instructions}} />
              </ListItem>
              ))
              }
             </List>
             </Paper>
         }
          <StreetView location={steps.steps[currentIndex].Loc} locEnd={steps.steps[currentIndex].EndLoc} />
        </Box>
        <Box alignItems="center" justifyContent="center" display="flex" >
                <Button className={classes.np_button} onClick={() => getPreviousStep()} variant="contained" color="primary" fullWidth>
                Previous Step
                </Button>
                <Button className={classes.np_button} onClick={() => getNextStep()} variant="contained" color="primary" fullWidth>
                Next Step
                </Button>
          </Box>
  	</div>
  );

}

export default NewDirectionsViewer