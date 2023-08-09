import React, { useState, useRef, useEffect } from 'react';
import clsx from 'clsx';
import { makeStyles } from '@material-ui/core/styles';
import CssBaseline from '@material-ui/core/CssBaseline';
import Box from '@material-ui/core/Box';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';
import Link from '@material-ui/core/Link';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import Button from '@material-ui/core/Button';
import Radio from '@material-ui/core/Radio';
import RadioGroup from '@material-ui/core/RadioGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormLabel from '@material-ui/core/FormLabel';
import DateFnsUtils from '@date-io/date-fns';
import ReactGA from 'react-ga';
import {
  MuiPickersUtilsProvider,
  KeyboardDatePicker,
} from '@material-ui/pickers';
import { isMobile } from 'react-device-detect';
import TextareaAutosize from '@material-ui/core/TextareaAutosize';
import { CSVLink } from "react-csv";
import CircularProgress from '@material-ui/core/CircularProgress';
import config from "./config.json"
import logo from "./logo192.png"
import MultiSelect from './MultiSelect';
const axios = require('axios');


function Copyright() {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {'Copyright © '}
      <Link color="inherit" href="#">
        BhavCopy Downloader
      </Link>{' '}
      {new Date().getFullYear()}
      {'.'}
    </Typography>
  );
}

const drawerWidth = 240;

const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
    justifyContent: "center",
  },
  toolbar: {
    paddingRight: 24, // keep right padding when drawer closed
    display: isMobile ? 'inline' : 'flex',
  },
  toolbarIcon: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
    padding: '0 8px',
    ...theme.mixins.toolbar,
  },
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
  },
  appBarShift: {
    marginLeft: drawerWidth,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  menuButton: {
    marginRight: 36,
  },
  menuButtonHidden: {
    display: 'none',
  },
  title: {
    flexGrow: 1,
  },
  drawerPaper: {
    position: 'relative',
    whiteSpace: 'nowrap',
    width: drawerWidth,
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  drawerPaperClose: {
    overflowX: 'hidden',
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    width: theme.spacing(7),
    [theme.breakpoints.up('sm')]: {
      width: theme.spacing(9),
    },
  },
  appBarSpacer: theme.mixins.toolbar,
  content: {
    height: '100vh',
  },
  container: {
    paddingTop: theme.spacing(4),
    paddingBottom: theme.spacing(4),
  },
  paper: {
    padding: theme.spacing(2),
    display: 'flex',
    overflow: 'auto',
    flexDirection: 'column',
  },
  fixedHeight: {
    height: 240,
  },
}));

export default function Dashboard() {
  useEffect(() => {
    ReactGA.pageview(window.location.pathname + window.location.search);
  }, []);
  const [showProgress, setShowProgress] = useState(false)
  const csvLink = useRef()
  const classes = useStyles();
  var date = new Date();
  date.setDate(date.getDate() - 1);
  const [selectedDate, setSelectedDate] = useState(date);
  const [index, setIndex] = useState("All");
  const [csvResponse, setCsvResponse] = useState([]);
  const [indexData, setIndexData] = useState([]);
  const handleDateChange = (date) => {
    setSelectedDate(date);
  };
  function handleIndexChange(e) {
    const indexName = e.target.value
    setIndex(indexName)
    if (indexName === "All") {
      setIndexData([]);
      return;
    }
    if (exchange === "bse" || exchange === "") {
      import('../BSE-Index-Configs/' + indexName).then((data) => {
        setIndexData(data.default);
      });
    } else {
      import('../NSE-Index-Configs/' + indexName).then((data) => {
        setIndexData(data.default);
      });
    }

  }
  const [exchange, setExchange] = React.useState('nse');
  function getDateInFormat() {
    return selectedDate.toLocaleDateString('en-GB', {
      day: '2-digit', month: 'short', year: 'numeric'
    }).replace(/ /g, '-')
  }
  const handleRadioChange = (event) => {
    const exchange = event.target.value;
    setExchange(exchange);
    setIndex("All")
    setIndexData([])
    if (exchange === "bse") {
      setFund(config.bseFund[0])
      return
    }
    setFund(config.nseFund[0])
  };
  function handleTextChange(e) {
    setIndexData(e.target.value.split(","))
  }
  function getDate() {
    let monthNames = ["Jan", "Feb", "Mar", "Apr",
      "May", "Jun", "Jul", "Aug",
      "Sep", "Oct", "Nov", "Dec"];

    let day = selectedDate.getDate();
    if (day < 10) {
      day = "0" + day
    }
    let monthIndex = selectedDate.getMonth();
    let monthName = monthNames[monthIndex];

    let year = selectedDate.getFullYear();
    return `${day}${monthName}${year}`;
  }
  function handleDownloadClick() {
    setShowProgress(true)
    const data = {
      "Date": getDate(),
      "Stocks": indexData,
      "Exchange": exchange.toUpperCase(),
      "Fund": fund
    }
    ReactGA.event({
      category: 'User',
      action: 'Download clicked with data :' + JSON.stringify(data)
    });
    if (fund !== "OPTIONS") {
      axios({
        method: 'post',
        url: config.backendUrl + '/getbhavcopy',
        data: data
      }).then(function (response) {
        setCsvResponse(response.data);
        csvLink.current.link.click();
        setShowProgress(false)
      }).catch(function (error) {
        console.log(error);
        setShowProgress(false)
      })
    } else {
      axios({
        method: 'post',
        url: config.backendUrl + '/optionChain?symbol=BANKNIFTY',
        data: data
      }).then(function (response) {
        exportToJson(response.data, "OptionChain_" + data["Date"] + ".json");
        setShowProgress(false)
      }).catch(function (error) {
        console.log(error);
        setShowProgress(false)
      })
    }

  }
  const [fund, setFund] = useState(config.nseFund[0])
  function handlefundChange(e) {
    setFund(e.target.value)
  }
  function handleSelect(values) {
    let selectedStocks = []
    values.forEach((value) => {
      selectedStocks.push(value.name.trim());
    })
    setIndexData(selectedStocks)
  }
  const fileName = getDateInFormat();
  return (
    <div className={classes.root}>
      <CssBaseline />
      <AppBar position="absolute" className={clsx(classes.appBar && classes.appBarShift)}>
        <Toolbar className={classes.toolbar}>
          <Typography component="h1" variant="h6" color="inherit" noWrap className={classes.title}>

            <img
              alt="bhavcopy downloader"
              href="#"
              style={{ height: "35px" }}
              src={logo}
            /><span style={{
              margin: "5px",
              verticalAlign: "top"
            }} >BhavCopy Downloader</span>
          </Typography>
          <span>
            <span> <a className="github-button" href="https://github.com/girishg4t/bhavCopy-downloader" data-size="large" aria-label="View girishg4t/bhavCopy-downloader on GitHub">View Source</a> {' '}</span>
            <span> <a className="github-button" href="https://github.com/girishg4t/bhavCopy-downloader" data-icon="octicon-star" data-size="large" data-show-count="true" aria-label="Star girishg4t/bhavCopy-downloader on GitHub">Star</a>{' '}</span>
            <a className="github-button" href="https://github.com/girishg4t/bhavCopy-downloader/fork" data-size="large" data-show-count="true" aria-label="Fork girishg4t/bhavCopy-downloader on GitHub">Fork</a>
          </span>
        </Toolbar>
      </AppBar>
      <main className={classes.content}>
        <div className={classes.appBarSpacer} />
        <Grid container spacing={1} style={{ margin: "20px" }}>
          <Grid container spacing={1} style={isMobile ? { flexWrap: "unset", display: "inline" } :
            { flexWrap: "unset" }}>
            <Grid item xs={3} style={{ maxWidth: "200px", paddingTop: "10px", flexWrap: "unset", alignSelf: "center" }}>
              <FormLabel component="legend">Stock Exchange</FormLabel>
              <RadioGroup aria-label="exchange" style={{ flexDirection: "inherit", flexWrap: "nowrap" }}
                name="exchange" value={exchange} onChange={handleRadioChange}>
                <FormControlLabel value="nse" control={<Radio />} label="NSE" />
                <FormControlLabel value="bse" control={<Radio />} label="BSE" />
              </RadioGroup>
            </Grid>
            <Grid item xs={3} style={{ alignSelf: "center" }}>
              <FormControl variant="outlined" className={classes.formControl}>
                <InputLabel id="demo-simple-select-outlined-label">Select Fund</InputLabel>
                <Select
                  labelId="demo-simple-select-outlined-label"
                  id="demo-simple-select-outlined"
                  value={fund}
                  defaultValue={fund}
                  onChange={handlefundChange}
                  label="Select Fund"
                  style={{ width: "250px" }}
                >
                  {
                    exchange === "nse" ? config.nseFund.map((fund) => {
                      return (<MenuItem value={fund}>{fund}</MenuItem>)
                    }) :
                      config.bseFund.map((fund) => {
                        return (<MenuItem value={fund}>{fund}</MenuItem>)
                      })
                  }

                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={3} style={{ alignSelf: "center" }}>
              <FormControl variant="outlined" className={classes.formControl}>
                <InputLabel id="demo-simple-select-outlined-label">Select Index</InputLabel>
                <Select
                  labelId="demo-simple-select-outlined-label"
                  id="demo-simple-select-outlined"
                  value={index}
                  defaultValue=""
                  onChange={handleIndexChange}
                  label="Select Indices"
                  style={{ width: "250px", height: "32px" }}
                >
                  <MenuItem value="All">
                    <em>All</em>
                  </MenuItem>
                  {
                    exchange === "nse" ? config.nseIndexs.map((index) => {
                      return (<MenuItem value={index + ".json"}>{index.replace(/_/g, " ")}</MenuItem>)
                    }) :
                      config.bseIndex.map((index) => {
                        return (<MenuItem value={index + ".json"}>{index.replace(/_/g, " ")}</MenuItem>)
                      })
                  }

                </Select>
              </FormControl>
              <MultiSelect handleSelect={handleSelect} />
            </Grid>
            <Grid item xs={3} style={{ alignSelf: "center" }}>
              <FormControl variant="outlined" className={classes.formControl} style={{ width: "250px" }}>
                <MuiPickersUtilsProvider utils={DateFnsUtils}>
                  <Grid container>
                    <KeyboardDatePicker
                      disableToolbar
                      variant="inline"
                      format="MM/dd/yyyy"
                      margin="normal"
                      id="date-picker-inline"
                      label="Date"
                      value={selectedDate}
                      onChange={handleDateChange}
                      KeyboardButtonProps={{
                        'aria-label': 'change date',
                      }}
                    />
                  </Grid>
                </MuiPickersUtilsProvider>
              </FormControl>
            </Grid>
            <Grid item xs={3} style={{ alignSelf: "center" }}>
              <Button disabled={showProgress} variant="contained" color="primary" onClick={handleDownloadClick}>
                Download
              </Button>
              {showProgress ? <CircularProgress /> : <div />}
              <CSVLink
                data={csvResponse}
                filename={index ? exchange + "-" + index.split(".")[0] + "-" + getDateInFormat() + ".csv" : exchange + "-" + fileName + ".csv"}
                className="btn btn-primary"
                ref={csvLink}
                target="_blank" />
            </Grid>
          </Grid>
        </Grid>
        <Grid spacing={1}>
          <TextareaAutosize style={{ width: "100%", height: "360px", fontSize: "large", overflow: "none" }} aria-label="maximum height"
            value={indexData}
            onChange={handleTextChange}
            placeholder="" />
        </Grid>
        <Box pt={4}>
          <Copyright />
        </Box>
      </main>
    </div>
  );
}

function exportToJson(objectData, filename) {
  let contentType = "application/json;charset=utf-8;";
  if (window.navigator && window.navigator.msSaveOrOpenBlob) {
    var blob = new Blob([decodeURIComponent(encodeURI(JSON.stringify(objectData)))], { type: contentType });
    navigator.msSaveOrOpenBlob(blob, filename);
  } else {
    var a = document.createElement('a');
    a.download = filename;
    a.href = 'data:' + contentType + ',' + encodeURIComponent(JSON.stringify(objectData));
    a.target = '_blank';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
  }
}