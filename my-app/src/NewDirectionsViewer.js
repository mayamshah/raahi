import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { GoogleMap, StreetViewPanorama} from '@react-google-maps/api'
import ListItem from '@material-ui/core/ListItem';
import Button from '@material-ui/core/Button';


function StreetView(props) {
	const mapContainerStyle = {
	  height: "400px",
	  width: "100%"
	};

	const center = {
	  lat:  props.location[0],
	  lng: props.location[1]
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
        <StreetView location={steps.steps[currentIndex].Loc}/>
  	</div>
  );

}

export default NewDirectionsViewer