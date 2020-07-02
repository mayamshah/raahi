import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import axios from "axios";

let endpoint = "http://localhost:8080/api/execute";

const useStyles = makeStyles((theme) => ({
  paper: {
    alignItems: 'center',
  },
  heading: {
    margin: theme.spacing(3,0,0,3),
    color: 'black',
    height: 48,
  },
  textField: {
    width: '25ch',
    margin: theme.spacing(0,0,3,3),
  },
  button: {
    margin: theme.spacing(3,0,0,3),
  },
}));

function Display() {
  const classes = useStyles();
  const [result, setResult] = useState("nothing");
  const [currentInput, setCurrentInput] = useState("");


  function getResult() {
     const request = {
        address: currentInput,
        distance: "1"
      };

      axios.post(endpoint, request)
      .then(res => {
        console.log(res);
        console.log(res.data);
        setResult(res.data.Error)
      })
  }

  function handleChange(e) {
    setCurrentInput(e.target.value)
  }

  return (
    <div>
        <form className={classes.textField} noValidate>
          <TextField id="standard-basic" label="Address" onChange={e => handleChange(e)}/>
        </form>
        <Button className={classes.button} onClick={() => getResult()} variant="contained" color="primary" >
        Enter
        </Button>  
        <Typography className={classes.heading} component="h1" variant="h6" gutterBottom>
           {result}
        </Typography>
    </div>
    );
}


function App() {

  const classes = useStyles();
  return (
    <div className={classes.paper}> 
      <Typography className={classes.heading} component="h1" variant="h6" gutterBottom>
      	What is your starting point?
      </Typography>
      <Display />
    </div>
  );
}


export default App;
