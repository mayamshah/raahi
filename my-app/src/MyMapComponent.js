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

// /*global google*/
// import React, { useState, useEffect } from 'react'
// import  { compose, withProps, lifecycle } from 'recompose'
// import {withScriptjs, withGoogleMap, GoogleMap, DirectionsRenderer} from 'react-google-maps'

// function createWayPoints(org) {

//     var waypoints = []
//     var i = 0
//     for (i = 0; i < org.length; i = i + 2) {
//       waypoints.push({location: new google.maps.LatLng(org[i], org[i+1]), stopover: true})
//     }

//     return waypoints
// }


// // function MyMapComponent(props) {
// //   const [currentProps, setCurrentProps] = useState(props.org);

// //   useEffect(() => { 
// //           const DirectionsService = new google.maps.DirectionsService();
// //           DirectionsService.route({
// //             origin: new google.maps.LatLng(currentProps.org[0], currentProps.org[1]),
// //             destination: new google.maps.LatLng(currentProps.org[0], currentProps.org[1]),
// //             waypoints: createWayPoints(this.props.org.slice(2)),
// //             travelMode: google.maps.TravelMode.WALKING,
// //           }, (result, status) => {
// //             if (status === google.maps.DirectionsStatus.OK) {
// //                 this.setState({
// //                   directions: {...result},
// //                   markers: true
// //                 })
// //             } else {
// //               console.error(`error fetching directions ${result}`);
// //             }
// //           });
// //         })

// //   function DirectionsComponent() {
// //      return compose(
// //       withProps({
// //         googleMapURL: "https://maps.googleapis.com/maps/api/js?key=AIzaSyB32cCcL4gD_WIYPP6dAVSprY_QYE3arsk",
// //         loadingElement: <div style={{ height: `400px` }} />,
// //         containerElement: <div style={{ width: `100%` }} />,
// //         mapElement: <div style={{height: `600px`, width: `600px` }}  />,
// //         org: this.props.org
// //       }),
// //       withScriptjs,
// //       withGoogleMap,

// //   }

// //   return (
// //     <DirectionsComponent
// //         />
// //   );
// // }

// const DirectionsComponent = compose(
//       withProps({
//         googleMapURL: "https://maps.googleapis.com/maps/api/js?key=AIzaSyB32cCcL4gD_WIYPP6dAVSprY_QYE3arsk",
//         loadingElement: <div style={{ height: `400px` }} />,
//         containerElement: <div style={{ width: `100%` }} />,
//         mapElement: <div style={{height: `600px`, width: `600px` }}  />,
//         // org: this.props.org
//       }),
//       withScriptjs,
//       withGoogleMap,
//       (props =>
//       <GoogleMap
//         defaultZoom={3}
//       >
//         {props.directions && <DirectionsRenderer directions={props.directions} suppressMarkers={props.markers}/>}
//       </GoogleMap>
//     ))

// class MyMapComponent extends React.PureComponent {

//   constructor(props) {
//     super(props); 
//   }
  

//   componentDidMount() {
//   new google.maps.DirectionsService().route({
//             origin: new google.maps.LatLng(this.props.org[0], this.props.org[1]),
//             destination: new google.maps.LatLng(this.props.org[0], this.props.org[1]),
//             waypoints: createWayPoints(this.props.org.slice(2)),
//             travelMode: google.maps.TravelMode.WALKING,
//           }, (result, status) => {
//             if (status === google.maps.DirectionsStatus.OK) {
//                 this.setState({
//                   directions: {...result},
//                   markers: true
//                 })
//                 console.log("inside")
//             } else {
//               console.error(`error fetching directions ${result}`);
//             }
//       })
//   }
  
//    render() { 
//     return (
//       <DirectionsComponent 
//        directions={this.state.directions} markers={this.state.markers}
//       />
//     )
//    }
// }

    

// export default MyMapComponent