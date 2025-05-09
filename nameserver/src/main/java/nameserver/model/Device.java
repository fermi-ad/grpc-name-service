package nameserver.model;

import io.quarkus.hibernate.orm.panache.PanacheEntity;
import jakarta.persistence.CascadeType;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.OneToMany;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;

@Entity
public class Device extends PanacheEntity {
    @NotBlank
    @Column(unique = true)
    private String name;

    private String description;
    
    @ManyToOne(optional = false)
    @JoinColumn(name = "node_id")
    private Node node;

    /*@OneToMany(mappedBy = "device", cascade = CascadeType.ALL, orphanRemoval = true)
    public java.util.List<Channel> channels;
*/
    public String getName() {
        return name;
    }

    public Device setName(String name) {
        this.name = name;
        return this;        
    }

    public String getDescription() {
        return description;
    }

    public Device setDescription(String description) {
        this.description = description;
        return this;        
    }

    public Node getNode() {
        return node;
    }

    public Device setNode(Node node) {
        this.node = node;
        return this;        
    }

    /*public java.util.List<Channel> getChannels() {
        return channels;
    }*/
    
    public Device() {
        // Default constructor
    }

    public Device(String name, String description, Node node) {
        this.name = name;
        this.description = description;
        this.node = node;
    }
}
