package nameserver.model;

import io.quarkus.hibernate.orm.panache.PanacheEntity;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.validation.constraints.NotBlank;
import jakarta.persistence.Table;

@Entity
public class ChannelTransform extends PanacheEntity {
    @NotBlank
    @Column(unique = true)
    public String name;

    @NotBlank    
    public String transform;
    
    @ManyToOne(optional = false)
    @JoinColumn(name = "channel_id")
    public Channel channel;
    
    public String description;

    public ChannelTransform() {
        // Default constructor
    }

    public ChannelTransform(String name, String transform, Channel channel) {
        this.name = name;
        this.transform = transform;
        this.channel = channel;
    }

    public String getName() {
        return name;
    }

    public ChannelTransform setName(String name) {
        this.name = name;
        return this;        
    }

    public String getTransform() {
        return transform;
    }

    public ChannelTransform setTransform(String transform) {
        this.transform = transform;
        return this;        
    }

    public Channel getChannel() {
        return channel;
    }

    public ChannelTransform setChannel(Channel channel) {
        this.channel = channel;
        return this;        
    }

    public String getDescription() {
        return description;
    }

    public ChannelTransform setDescription(String description) {
        this.description = description;
        return this;        
    }
}

