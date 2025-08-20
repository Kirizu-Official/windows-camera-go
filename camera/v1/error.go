package camera

func (d *Device) SetError(err error) {
	if d.errorCallBack != nil {
		d.errorCallBack(err)
	}
}
