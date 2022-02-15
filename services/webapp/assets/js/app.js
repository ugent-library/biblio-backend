import htmx from 'htmx.org';
import csrf from './ui/csrf.js'
import checkbox from './ui/checkbox.js'
import bootstrap from './ui/bootstrap.js'
import popover from './ui/popover.js'
import header from './ui/header.js'
import multiple from './ui/multiple.js'
import changeSubmit from './ui/form_change_submit.js'
import autocomplete from './ui/autocomplete.js';
import submit from './ui/form_submit.js'
import modalClose from './ui/modal_close.js'
import modalPopper from './ui/modal_popper.js'
import tabs from './ui/tabs.js'
import radioCard from './ui/radio_card.js'
import membership from './ui/membership.js'
import toast from './ui/toast.js'
import fileEmbargo from './ui/file_embargo.js'
import datasetEmbargo from './ui/dataset_embargo.js';
import fileShowCCLicense from './ui/file_show_cc_license.js';

document.addEventListener('DOMContentLoaded', function () {
    csrf()
    tabs()
    checkbox()
    bootstrap()
    popover()
    header()
    multiple()
    changeSubmit()
    autocomplete()
    submit()
    modalClose()
    modalPopper()
    radioCard()
    membership()
    toast()
    fileEmbargo()
    datasetEmbargo()
    fileShowCCLicense()
});