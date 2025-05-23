package nameserver.model;

import io.quarkus.hibernate.orm.panache.PanacheEntity;
import jakarta.persistence.Entity;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.JoinTable;
import jakarta.persistence.ManyToMany;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.Table;

@Entity
public class ChannelAccessControl extends PanacheEntity {
    @ManyToOne(optional = false)
    @JoinColumn(name = "channel_id")
    private Channel channel;
    
    @ManyToOne(optional = false)
    @JoinColumn(name = "role_id")
    private Role role;
    
    private boolean read;
    private boolean write;

    public ChannelAccessControl() {
        // Default constructor
    }
    
    public Channel getChannel() {
        return channel;
    }

    public ChannelAccessControl setChannel(Channel channel) {
        this.channel = channel;
        return this;        
    }

    public Role getRole() {
        return role;
    }

    public ChannelAccessControl setRole(Role role) {
        this.role = role;
        return this;        
    }  

    public boolean isRead() {
        return read;
    }

    public ChannelAccessControl setRead(boolean read) {
        this.read = read;
        return this;        
    }

    public boolean isWrite() {
        return write;
    }

    public ChannelAccessControl setWrite(boolean write) {
        this.write = write;
        return this;        
    }
}
