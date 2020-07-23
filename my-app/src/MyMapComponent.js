/*global google*/
import React, {useState} from 'react';
import { makeStyles } from '@material-ui/core/styles';
import  { compose, withProps, lifecycle } from 'recompose'
import {withScriptjs, withGoogleMap, GoogleMap, DirectionsRenderer} from 'react-google-maps'
import Typography from '@material-ui/core/Typography';
import PropTypes from 'prop-types';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Paper from '@material-ui/core/Paper';
import CssBaseline from '@material-ui/core/CssBaseline';
import Button from '@material-ui/core/Button';




function createWayPoints(org) {

    var waypoints = []
    var i = 0
    for (i = 0; i < org.length && i < 46; i = i + 2) {
      waypoints.push({location: new google.maps.LatLng(org[i], org[i+1]), stopover: false})
    }

    return waypoints
}

class MyMapComponent extends React.Component {
  constructor(props){
    super(props)
  }

  
render() {
    const DirectionsComponent = compose(
      withProps({
        googleMapURL: "https://maps.googleapis.com/maps/api/js?key=AIzaSyB32cCcL4gD_WIYPP6dAVSprY_QYE3arsk",
        loadingElement: <div style={{ height: `100%` }} />,
        containerElement: <div style={{ height: `100%`}} />,
        mapElement: <div style={{ height: `500px`, width: `100%`}}  />,
        response: this.props.response,
        steps: this.props.response.Directions
      }),
      withScriptjs,
      withGoogleMap,
      lifecycle({
        componentDidMount() { 
          console.log(this.props.response)
          const DirectionsService = new google.maps.DirectionsService();
          DirectionsService.route({
            origin: new google.maps.LatLng(this.props.response.Org[0], this.props.response.Org[1]),
            destination: new google.maps.LatLng(this.props.response.Dest[0], this.props.response.Dest[1]),
            waypoints: createWayPoints(this.props.response.Path),
            travelMode: google.maps.TravelMode.WALKING,
          }, (result, status) => {
            if (status === google.maps.DirectionsStatus.OK) {
            this.setState({
                directions: {...result},
                markers: true,
              })
            } else {
              console.error(`error fetching directions ${result}`);
            }
          });
        }
      })
    )(props =>
    
     <React.Fragment>
      <CssBaseline />
      {props.directions && 
        <div>
        <GoogleMap
          defaultZoom={3}
        >
          <DirectionsRenderer directions={props.directions} suppressMarkers={props.markers}/>
        </GoogleMap>
        </div>
        }
        
      </React.Fragment>
    );

return (
        <div>
        <DirectionsComponent/> 
        </div>
    )
  }
}



export default MyMapComponent

