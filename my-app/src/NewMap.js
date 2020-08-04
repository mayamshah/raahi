import React, { useState } from 'react'
import { GoogleMap, DirectionsService, DirectionsRenderer } from '@react-google-maps/api'


function LatLng(lat, lng) {
  return {lat: lat, lng: lng};
}
function createWayPoints(org) {

    var waypoints = []
    var i = 0
    for (i = 0; i < org.length && i < 46; i = i + 2) {
      waypoints.push({location: LatLng(org[i], org[i+1]), stopover: false})
    }

    return waypoints
}

function DirectionsNew(props) {
    
    const [current_response, setResponse] = useState({});

    // const state = React.memo(props => {

    //   return  {
    //             response: null,
    //             travelMode: 'WALKING',
    //             origin: "12543 Palmtag Drive, Saratoga, CA",
    //             destination: "12546 Palmtag Drive, Saratoga, CA",
    //           };
    //   });

    const state = {
                response: null,
                travelMode: 'WALKING',
                origin: LatLng(props.response.Org[0], props.response.Org[1]),
                destination: LatLng(props.response.Dest[0], props.response.Dest[1]),
                waypoints: createWayPoints(props.response.Path)
              };

  // shouldComponentUpdate(nextProps, nextState){
  //   if (nextProps.response.Org[0] === this.props.response.Org[0] || this.state.response != null || this.state.response == nextState.response) {
  //     console.log("entered");
  //     return false 
  //   }
  //   return true
  // }

  function directionsCallback(response) {
    console.log("callback")
    console.log(response)

    if (response !== null && current_response.status == null) {
      if (response.status === 'OK') {
        console.log("inside")
        setResponse(response)
      } else {
        console.log('response: ', response)
      }
    }
  }

  function onMapClick(...args) {
    console.log('onClick args: ', args)
  }

    return (
      console.log("return"),
      <div>
      <div className='map'>
        <div className='map-settings'>
          <hr className='mt-0 mb-3' />

        </div>

        <div className='map-container'>
          <GoogleMap
            // required
            id='direction-example'
            // required
            mapContainerStyle={{
              height: '400px',
              width: '100%'
            }}
            // required
            zoom={2}
            // required
            center={{
              lat: 0,
              lng: -180
            }}
            // optional
            onClick={onMapClick}
            // optional
            onLoad={map => {
              console.log('DirectionsRenderer onLoad map: ', map)
            }}
            // optional
            onUnmount={map => {
              console.log('DirectionsRenderer onUnmount map: ', map)
            }}
          >
            
                <DirectionsService
                  // required
                  options={{ 
                    destination: state.destination,
                    origin: state.origin,
                    travelMode: state.travelMode,
                    waypoints: state.waypoints
                  }}
                  // required
                  callback={directionsCallback}
                  // optional
                  onLoad={directionsService => {
                    console.log(state);
                    console.log('DirectionsService onLoad directionsService: ', directionsService)
                  }}
                  // optional
                  onUnmount={directionsService => {
                    console.log('DirectionsService onUnmount directionsService: ', directionsService)
                  }}
                />

            {
              current_response !== null && (
                <DirectionsRenderer
                  // required
                  options={{ 
                    directions: current_response
                  }}
                  // optional
                  onLoad={directionsRenderer => {
                    console.log('DirectionsRenderer onLoad directionsRenderer: ', directionsRenderer)

                  }}
                  // optional
                  onUnmount={directionsRenderer => {
                    console.log('DirectionsRenderer onUnmount directionsRenderer: ', directionsRenderer)
                  }}
                />
              )
            }
          </GoogleMap>
        </div>
      </div>
    </div>
    )
  
}



export default DirectionsNew