/*global google*/
import React from 'react'
import  { compose, withProps, lifecycle } from 'recompose'
import {withScriptjs, withGoogleMap, GoogleMap, DirectionsRenderer} from 'react-google-maps'

function createWayPoints(org) {

    var waypoints = []
    var i = 0
    for (i = 0; i < org.length; i = i + 2) {
      waypoints.push({location: new google.maps.LatLng(org[i], org[i+1]), stopover: true})
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
        loadingElement: <div style={{ height: `400px` }} />,
        containerElement: <div style={{ width: `100%` }} />,
        mapElement: <div style={{height: `600px`, width: `600px` }}  />,
        org: this.props.org
      }),
      withScriptjs,
      withGoogleMap,
      lifecycle({
        componentDidMount() { 
          const DirectionsService = new google.maps.DirectionsService();
          DirectionsService.route({
            origin: new google.maps.LatLng(this.props.org[0], this.props.org[1]),
            destination: new google.maps.LatLng(this.props.org[0], this.props.org[1]),
            waypoints: createWayPoints(this.props.org.slice(2)),
            travelMode: google.maps.TravelMode.WALKING,
          }, (result, status) => {
            if (status === google.maps.DirectionsStatus.OK) {
this.setState({
                directions: {...result},
                markers: true
              })
            } else {
              console.error(`error fetching directions ${result}`);
            }
          });
        }
      })
    )(props =>
      <GoogleMap
        defaultZoom={3}
      >
        {props.directions && <DirectionsRenderer directions={props.directions} suppressMarkers={props.markers}/>}
      </GoogleMap>
    );
return (
        <DirectionsComponent
        />
    )
  }
}
export default MyMapComponent