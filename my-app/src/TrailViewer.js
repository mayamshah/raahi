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

import Typography from '@material-ui/core/Typography';
import Avatar from '@material-ui/core/Avatar';
import { FixedSizeList } from 'react-window';





function TrailViewer(trails) {
	console.log(trails)

	const [currentIndex, setCurrentIndex] = useState(0);

	function changeStep(index){
      setCurrentIndex(index)
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