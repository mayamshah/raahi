import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { GoogleMap} from '@react-google-maps/api'
import ListItem from '@material-ui/core/ListItem';
import Button from '@material-ui/core/Button';
import { Box } from '@material-ui/core';
import Paper from '@material-ui/core/Paper';
import List from '@material-ui/core/List';
import ListItemAvatar from '@material-ui/core/ListItemAvatar';
import ListItemText from '@material-ui/core/ListItemText';
import Divider from '@material-ui/core/Divider';
import getDirections from './DirectionsExporter.js'
import Typography from '@material-ui/core/Typography';
import Avatar from '@material-ui/core/Avatar';
import { FixedSizeList } from 'react-window';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import IconButton from '@material-ui/core/IconButton';
import DirectionsIcon from '@material-ui/icons/Directions';
import Tooltip from '@material-ui/core/Tooltip';


function TrailViewer(trails) {
	console.log(trails)

	const [currentIndex, setCurrentIndex] = useState(0);

	function changeStep(index){
      setCurrentIndex(index)
    }

    function exportDirections() {
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
          value: "driving"        // may be "walking", "bicycling" or "transit" as well
        },
        {
          key: "dir_action",
          value: "navigate"       // this instantly initializes navigation using the given travel mode
        }
      ], 
      }
    data.source.latitude = trails.origin[0]
    data.source.longitude =  trails.origin[1]
    data.destination.latitude = trails.trails[currentIndex].Coords[0]
    data.destination.longitude = trails.trails[currentIndex].Coords[1]
   
    getDirections(data)
  }

	return (
		<Box alignItems="center" justifyContent="center" display="flex" >
          {trails.trails && 
            <Paper style={{maxHeight: 400, overflow: 'auto', height:400, width: 400}}>
            <List>
             {trails.trails.map((trail, index) => (
              	<ListItem button selected={index === currentIndex} onClick={() => changeStep(index)}>
              	<ListItemAvatar>
          			<Avatar style={{fontSize: '15px'}} variant="rounded">
          			{Number(trail.DistFromOrg).toFixed(2).toString()}
          			</Avatar>
       			</ListItemAvatar>
				<ListItemText primary={trail.Name} secondary={trail.Summary} />
				<Divider variant="inset" component="li" />
				{index === currentIndex &&	
					<ListItemIcon>
					<Tooltip title="Directions">
						<IconButton onClick={() => exportDirections()}>
	                        <DirectionsIcon fontSize="large"/>
	                    </IconButton>
	                 </Tooltip>
					</ListItemIcon>}
				</ListItem>
              ))
             }
             </List>
             </Paper>
         }
         <GoogleMap
		    id="circle-example"
		    mapContainerStyle={{
		      height: "400px",
		      width: "400px"
		    }}
		    zoom={15}
		    center={{
		      lat: trails.trails[currentIndex].Coords[0],
		      lng: trails.trails[currentIndex].Coords[1]
		    }}
		  />
        </Box>
	)
}

export default TrailViewer