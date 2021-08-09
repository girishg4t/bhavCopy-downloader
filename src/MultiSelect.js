/* eslint-disable no-use-before-define */
import React from "react";
import { fade, makeStyles } from "@material-ui/core/styles";
import Popper from "@material-ui/core/Popper";
import SettingsIcon from "@material-ui/icons/Settings";
import CloseIcon from "@material-ui/icons/Close";
import DoneIcon from "@material-ui/icons/Done";
import Autocomplete from "@material-ui/lab/Autocomplete";
import ButtonBase from "@material-ui/core/ButtonBase";
import InputBase from "@material-ui/core/InputBase";
import stocks from './NSE-Stocks/stocks.js';

const useStyles = makeStyles((theme) => ({
    root: {
        paddingRight: 5,
        paddingLeft: 5,
        paddingTop: 5
    },
    button: {
        fontSize: 16,
        width: "100%",
        textAlign: "left",
        paddingBottom: 8,
        color: "blue",
        "&:hover,&:focus": {
            color: "#0366d6"
        },
        "& span": {
            width: "100%"
        },
        "& svg": {
            width: 18,
            height: 18
        }
    },
    popper: {
        border: "1px solid rgba(27,31,35,.15)",
        boxShadow: "0 3px 12px rgba(27,31,35,.15)",
        borderRadius: 3,
        width: 300,
        zIndex: 1,
        fontSize: 13,
        color: "#586069",
        backgroundColor: "#f6f8fa"
    },
    header: {
        borderBottom: "1px solid #e1e4e8",
        padding: "8px 10px",
        fontWeight: 600
    },
    inputBase: {
        padding: 10,
        width: "100%",
        borderBottom: "1px solid #dfe2e5",
        "& input": {
            borderRadius: 4,
            backgroundColor: theme.palette.common.white,
            padding: 8,
            transition: theme.transitions.create(["border-color", "box-shadow"]),
            border: "1px solid #ced4da",
            fontSize: 14,
            "&:focus": {
                boxShadow: `${fade(theme.palette.primary.main, 0.25)} 0 0 0 0.2rem`,
                borderColor: theme.palette.primary.main
            }
        }
    },
    paper: {
        boxShadow: "none",
        margin: 0,
        color: "#586069",
        fontSize: 13
    },
    option: {
        minHeight: "auto",
        alignItems: "flex-start",
        padding: 8,
        '&[aria-selected="true"]': {
            backgroundColor: "transparent"
        },
        '&[data-focus="true"]': {
            backgroundColor: theme.palette.action.hover
        }
    },
    popperDisablePortal: {
        position: "relative"
    },
    iconSelected: {
        width: 17,
        height: 17,
        marginRight: 5,
        marginLeft: -2
    },
    text: {
        flexGrow: 1
    },
    close: {
        opacity: 0.6,
        width: 18,
        height: 18
    }
}));

export default function MultiSelect({handleSelect}) {
    const classes = useStyles();
    const [anchorEl, setAnchorEl] = React.useState(null);
    const [value, setValue] = React.useState([]);
    const [pendingValue, setPendingValue] = React.useState([]);

    const handleClick = (event) => {
        setPendingValue(value);
        setAnchorEl(event.currentTarget);
    };

    const handleClose = (event, reason) => {
        if (reason === "toggleInput") {
            return;
        }
        setValue(pendingValue);
        if (anchorEl) {
            anchorEl.focus();
        }
        setAnchorEl(null);
        handleSelect(pendingValue);
        event.preventDefault();
    };

    const open = Boolean(anchorEl);
    const id = open ? "github-label" : undefined;

    return (
        <React.Fragment>
            <div className={classes.root}>
                <ButtonBase
                    disableRipple
                    className={classes.button}
                    aria-describedby={id}
                    onClick={handleClick}
                >
                    <div style={{paddingRight : "10px"}}>  or Search Stocks </div> <SettingsIcon />
                </ButtonBase>
            </div>
            <Popper
                id={id}
                open={open}
                anchorEl={anchorEl}
                placement="bottom-start"
                className={classes.popper}
            >
                <div className={classes.header}>Select stocks to be downloaded</div>
                <Autocomplete
                    open
                    onClose={handleClose}
                    multiple
                    classes={{
                        paper: classes.paper,
                        option: classes.option,
                        popperDisablePortal: classes.popperDisablePortal
                    }}
                    value={pendingValue}
                    onChange={(event, newValue) => {
                        setPendingValue(newValue);
                    }}
                    disableCloseOnSelect
                    disablePortal
                    renderTags={() => null}
                    noOptionsText="No labels"
                    renderOption={(option, { selected }) => (
                        <React.Fragment>
                            <DoneIcon
                                className={classes.iconSelected}
                                style={{ visibility: selected ? "visible" : "hidden" }}
                            />
                            <div className={classes.text}>
                                {option.name}
                                <br />
                                {option.description}
                            </div>
                            <CloseIcon
                                className={classes.close}
                                style={{ visibility: selected ? "visible" : "hidden" }}
                            />
                        </React.Fragment>
                    )}
                    options={[...labels].sort((a, b) => {
                        // Display the selected labels first.
                        let ai = value.indexOf(a);
                        ai = ai === -1 ? value.length + labels.indexOf(a) : ai;
                        let bi = value.indexOf(b);
                        bi = bi === -1 ? value.length + labels.indexOf(b) : bi;
                        return ai - bi;
                    })}
                    getOptionLabel={(option) => option.name}
                    renderInput={(params) => (
                        <InputBase
                            ref={params.InputProps.ref}
                            inputProps={params.inputProps}
                            autoFocus
                            className={classes.inputBase}
                        />
                    )}
                />
            </Popper>
        </React.Fragment>
    );
}

const labels = stocks();