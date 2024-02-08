package event

import "reflect"

type Plan struct {
	notifyItems []Item //notify关联
	eventItems  []Item //event关联
	Date        string
	UnitUFI     string
	PreFlightID string
	CallSign    string
	CallSigns   string
	RegNumber   string
}

func (p *Plan) MsgType() string {
	return "notify.plan"
}

func (p *Plan) MsgSubType() string {
	return "notify.plan.statusChanged"
}

func (p *Plan) PubBus() {
	Push()
}

func (p *Plan) Set(k, v string) {
	rv := reflect.ValueOf(p)
	rv.Elem().FieldByName(k).SetString(v)
}

func (p *Plan) Updates() error {

	//将当前内容发送消息到消息总线中
	p.PubBus()

	p.NotifyUp() //消息通知，同步更新关联
	p.EventUp()  //内部事件通知，同步更新关联事件

	return nil
}

func (p *Plan) NotifyUp() error {
	if len(p.notifyItems) == 0 {
		return nil
	}

	for _, v := range p.notifyItems {
		err := v.Updates()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Plan) EventUp() error {
	if len(p.eventItems) == 0 {
		return nil
	}

	for _, v := range p.eventItems {
		err := v.Updates()
		if err != nil {
			return err
		}
	}
	return nil
}

// AddItems 增加关联操作涉及
func (this *Plan) AddItems(item Item) {
	this.notifyItems = append(this.notifyItems, item)
}
