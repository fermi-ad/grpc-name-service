package nameserver.model;

import io.quarkus.hibernate.orm.panache.PanacheEntity;
import jakarta.persistence.CascadeType;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.OneToMany;
import jakarta.validation.constraints.NotBlank;

@Entity
public class Channel extends PanacheEntity {
    @NotBlank
    @Column(unique = true)
    private String name;
    private String description;
    private String metadata;

    @ManyToOne(optional = false)
    @JoinColumn(name = "device_id")
    private Device device;

    @OneToMany(mappedBy = "channel", cascade = CascadeType.ALL, orphanRemoval = true)
    private java.util.List<ChannelAlarm> alarms;

    @OneToMany(mappedBy = "channel", cascade = CascadeType.ALL, orphanRemoval = true)
    private java.util.List<ChannelTransform> transforms;

    @OneToMany(mappedBy = "channel", cascade = CascadeType.ALL, orphanRemoval = true)
    private java.util.List<ChannelAccessControl> accessControls;
    
    public String getName() {
        return name;
    }

    public Channel setName(String name) {
        this.name = name;
        return this;        
    }

    public String getDescription() {
        return description;
    }

    public Channel setDescription(String description) {
        this.description = description;
        return this;        
    }

    public String getMetadata() {
        return metadata;
    }

    public Channel setMetadata(String metadata) {
        this.metadata = metadata;
        return this;        
    }

    public Device getDevice() {
        return device;
    }

    public Channel setDevice(Device device) {
        this.device = device;
        return this;        
    }

    public java.util.List<ChannelAlarm> getAlarms() {
        return alarms;
    }


    public Channel addAlarm(ChannelAlarm alarm) {
        if (this.alarms == null) {
            this.alarms = new java.util.ArrayList<>();
        }
        this.alarms.add(alarm);
        return this;        
    }

    public Channel setAlarms(java.util.List<ChannelAlarm> alarms) {
        this.alarms = alarms;
        return this;        
    }
    
    public java.util.List<ChannelTransform> getTransforms() {
        return transforms;
    }

    public Channel addTransform(ChannelTransform transform) {
        if (this.transforms == null) {
            this.transforms = new java.util.ArrayList<>();
        }
        this.transforms.add(transform);
        return this;        
    }

    public Channel setTransforms(java.util.List<ChannelTransform> transforms) {
        this.transforms = transforms;
        return this;        
    }

    public java.util.List<ChannelAccessControl> getAccessControls() {
        return accessControls;
    }

    public Channel addAccessControl(ChannelAccessControl accessControl) {
        if (this.accessControls == null) {
            this.accessControls = new java.util.ArrayList<>();
        }
        this.accessControls.add(accessControl);
        return this;        
    }

    public Channel() {
        this.name = "";
        this.description = "";
        this.metadata = "";
    }

    public Channel(String name, String description, String metadata, Device device) {
        this.name = name;
        this.description = description;
        this.metadata = metadata;
        this.device = device;
    }
}
