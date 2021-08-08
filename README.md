# BhavCopy Downloader
Download NSE and BSE data free, Here is the [blog](https://girishg4t.github.io/2021/03/15/bhavcopy-downloader.html) post for the same

![logo](./public/favicon.ico)

BhavCopy downloader is a FREE NSE and BSE end of the day stocks data downloader. Since it connect to NSE and BSE server for getting the data it is considered as authentic. it can download daily as well as historical EOD data for the currently configured indexes. It is made simple, customizable, public and free, so any one can download the data as per there need. You can check the existing data in "./nse" and "./bse" folder which got downloaded so far.

### Feathers:
- Download EOD data for both the index NSE and BSE with delivered quantity and delivered quantity percentage
- Can select different funds like Equities, f&o etc.
- Can select different index from NSE and BSE like NIFTY50, BSE100 etc.
- Get back date data
- Api is public any one can use directly, https://bhavcopy-backend.herokuapp.com

### Flow:

![Alt working](./flow-diagram.png)
#### Backend: 
It is written in golang which makes an api call to NSE and BSE servers to get all stocks eod data. ones the data is received it is stored in this [github repo](https://github.com/girishg4t/nse-bse-bhavcopy). Further reading of data is happened from this files if present else it get downloaded and stored. Based on the api request, the data is then send in csv format.

#### Frontend:
It is written in javascript react, which allows user's to download the data based on there selection. Below are the steps user need to perform to get the data, default is all stock download. 
  
    
1) Select the Stock Exchange from which the data is required eg. NSE/BSE
2) Select Fund for the particular exchange, currently it is configured for Equity only, in future more options will be added
3) Select the Index for which the data is required default is All.
4) One's index is selected all the stock in that index get appear in textarea which is editable. User can make changes to the list.  
Current configured index are:  
<table>
  <tr>
    <td colspan="5">NSE</td>
    <td colspan="5">BSE</td>
  </tr>
  <tr>
    <td colspan="5">"AUTO",
        "BANK",
        "CONSUMERDURABLES",
        "FINANCE",
        "FINANCIAL_SERVICES",
        "FMCG",
        "HEALTHCARE",
        "IT",
        "MEDIA_ENTERTAINMENT",
        "METAL",
        "OIL_GAS",
        "PHARMA",
        "PRIVATE_BANK",
        "PSU_BANK",
        "REALTY",
        "NIFTY50",
        "NIFTY100",
        "NIFTY200",
        "NIFTY500",
        "NIFTY500_MULTICAP_50_25_25",
        "NIFTY_LARGE_MIDCAP250",
        "NIFTY_MIDCAP50",
        "NIFTY_MIDCAP100",
        "NIFTY_MIDCAP150",
        "NIFTY_MID_SMALLCAP400",
        "NIFTY_NEXT50",
        "NIFTY_SMALLCAP50",
        "NIFTY_SMALLCAP100",
        "NIFTY_SMALLCAP250"</td>
         <td colspan="5">"AUTO",
        "BANKS",
        "BASIC_MATERIALS",
        "CAPITAL_GOODS",
        "CONSUMER_DISCRETIONARY_GOODS_SERVICES",
        "CONSUMER_DURABLES",
        "ENERGY",
        "FINANCE",
        "FMCG",
        "HEALTHCARE",
        "INDUSTRIALS",
        "IT",
        "METAL",
        "OIL_GAS",
        "POWER",
        "REALTY",
        "TECK",
        "TELECOM",
        "UTILITIES",
        "BSE100",
        "BSE200",
        "BSE500",
        "BSE_ AllCap",
        "BSE_LARGECAP",
        "BSE_LARGE_MIDCAP",
        "BSE_MIDCAP",
        "BSE_MIDCAP_SELECT_INDEX",
        "BSE_MID_SMALLCAP",
        "BSE_SENSEX",
        "BSE_SENSEX_50",
        "BSE_SMALLCAP",
        "BSE_SMALLCAP_SELECT_INDEX"</td>
  </tr>
</table>

5) Select date, to download the specific day data, default is previous date 
6) Click on download to get the bhavcopy in csv file.

### In Action 

![Alt working](./working.gif)
### How to run

```sh
npm start
```

## Contributing

As you know stocks in an index get change and keeping that up to date is not a one person job, that's why i have made all this configurable.
There are 2 folders "./NSEIndexConfigs"  and "./BSEIndexConfigs" in which all the index stocks configuration are present eg. Auto, Nifty50, BSE100 etc.

You can make changes to config by adding/removing the stock as per NSE and BSE changes or create the new/customized config and raise the pull request.

Also there is "./src/config.json" in which UI related config are present.

Inputs are always welcome! whether it's.
- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features