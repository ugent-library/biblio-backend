import csrf from './ui/csrf.js'
import checkbox from './ui/checkbox.js'
import bootstrap from './ui/bootstrap.js'
import header from './ui/header.js'
import multiple from './ui/multiple.js'
import changeSubmit from './ui/form_change_submit.js'
import submit from './ui/form_submit.js'
import modalClose from './ui/modal_close.js'
import modalPopper from './ui/modal_popper.js'
import tabs from './ui/tabs.js'
import radioCard from './ui/radio_card.js'
import membership from './ui/membership.js'
import toast from './ui/toast.js'

document.addEventListener('DOMContentLoaded', function () {
    csrf()
    tabs()
    checkbox()
    bootstrap()
    header()
    multiple()
    changeSubmit()
    submit()
    modalClose()
    modalPopper()
    radioCard()
    membership()
    toast()
});