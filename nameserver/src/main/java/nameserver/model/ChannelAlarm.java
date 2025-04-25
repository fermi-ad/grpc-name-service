package nameserver.model;

import io.quarkus.hibernate.orm.panache.PanacheEntity;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.validation.constraints.NotBlank;
import jakarta.persistence.Table;
import jakarta.persistence.UniqueConstraint;

@Table(name = "channel_alarm", uniqueConstraints = { @UniqueConstraint(columnNames = {"alarm_type_id", "channel_id" }) })  
@Entity
public class ChannelAlarm extends PanacheEntity {
    @NotBlank
    @Column(name = "trigger_condition")
    public String triggerCondition;

    @ManyToOne(optional = false)
    @JoinColumn(name = "alarm_type_id")
    public AlarmType alarmType;
    
    @ManyToOne(optional = false)
    @JoinColumn(name = "channel_id")
    public Channel channel;
    
    public ChannelAlarm() {
        // Default constructor
    }
    public String getTriggerCondition() {
        return triggerCondition;
    }

    public ChannelAlarm setTriggerCondition(String triggerCondition) {
        this.triggerCondition = triggerCondition;
        return this;        
    }

    public AlarmType getAlarmType() {
        return alarmType;
    }

    public ChannelAlarm setAlarmType(AlarmType alarmType) {
        this.alarmType = alarmType;
        return this;        
    }

    public Channel getChannel() {
        return channel;
    }

    public ChannelAlarm setChannel(Channel channel) {
        this.channel = channel;
        return this;        
    }
}




    
