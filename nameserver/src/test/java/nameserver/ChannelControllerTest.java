package nameserver;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;

import java.util.ArrayList;
import java.util.List;

import org.junit.jupiter.api.Test;

import io.quarkus.test.TestTransaction;
import io.quarkus.test.junit.QuarkusTest;
import jakarta.inject.Inject;
import nameserver.model.Channel;
import nameserver.model.Device;
import nameserver.model.Node;

/*
 * Tests channel API using ChannelController.  Test environment prepopulates the database using 
 * import.sql.
 */
@QuarkusTest
public class ChannelControllerTest {

    @Inject
    ChannelController channelController; // Inject the ChannelController directly

    //Data from import.sql
    static final String TEST_DEVICE_NAME = "Device A";
    static final String TEST_CHANNEL_NAME = "A:ChannelA";
    static final String TEST_NODE_NAME = "Node A";

    private void createTestDevice(String name, List<String> channels) {
        Device device = new Device();
        device.setName(name);
        device.setDescription("Test device");
        Node node = Node.find("hostname", TEST_NODE_NAME).firstResult();
        assertNotNull(node);
        device.setNode(node);
        device.persist();

        for (String channelName : channels) {
            Channel channel = new Channel();
            channel.setName(channelName);
            channel.setDescription("Test channel");
            channel.setDevice(device);
            channel.persist();
        }
    }
    @Test
    @TestTransaction
    public void testListChannels() {       

        //add extra channels
        var channames = new ArrayList<String>();
        final var bchannelCount = 8;
        final var totalChanCount = 1 + bchannelCount;
        for(int i = 0; i < bchannelCount; i++) {
            channames.add("B:Channel" + i);
        }
        createTestDevice("DeviceB", channames);

        //test filtering by device name
        proto.ListChannelsRequest request = proto.ListChannelsRequest.newBuilder()
                .setDeviceName(TEST_DEVICE_NAME) 
                .build();

        var channels = ChannelController.listChannels(request).getChannelsList();

        assertNotNull(channels);
        assertEquals(1, channels.size());
        assertEquals(TEST_CHANNEL_NAME, channels.get(0).getName());
        assertEquals(TEST_DEVICE_NAME, channels.get(0).getDeviceName());

        //test with no filter
        request = proto.ListChannelsRequest.newBuilder()                
                .build();
        channels = ChannelController.listChannels(request).getChannelsList();
        assertNotNull(channels);
        assertEquals(channels.size(), totalChanCount);
        
        //test pagination
        request = proto.ListChannelsRequest.newBuilder()                
                .setPagination(proto.PaginationRequest.newBuilder()
                        .setPageSize(5)
                        .setPage(1))
                .build();

        var response = ChannelController.listChannels(request);
        channels = response.getChannelsList();
        var pagresp = response.getPagination();
        assertNotNull(channels);
        assertEquals(5, channels.size());
        assertEquals(1, pagresp.getPage()); 
        assertEquals(totalChanCount, pagresp.getTotalCount());
        
        request = proto.ListChannelsRequest.newBuilder()                
                .setPagination(proto.PaginationRequest.newBuilder()
                        .setPageSize(5)
                        .setPage(2))
                .build();
        response = ChannelController.listChannels(request);
        channels = response.getChannelsList();
        pagresp = response.getPagination();
        assertNotNull(channels);
        assertEquals(4, channels.size());
        assertEquals(2, pagresp.getPage());
        assertEquals(totalChanCount, pagresp.getTotalCount());

    }

    @Test
    @TestTransaction
    public void testCreateChannel() {
        // Create a new channel

        proto.Channel channel = proto.Channel.newBuilder()
                .setName("New Channel")
                .setDescription("A new test channel")                
                .setDeviceName(TEST_DEVICE_NAME) 
                .build();
        proto.CreateChannelRequest request = proto.CreateChannelRequest.newBuilder()
                .setChannel(channel)
                .build();
        
        var response = channelController.createChannel(request);

        assertNotNull(response.getChannel());
        var createdChannel = response.getChannel();
        assertEquals("New Channel", createdChannel.getName());
        assertEquals("A new test channel", createdChannel.getDescription());

        // Verify the channel was persisted
        Channel persistedChannel = Channel.find("name", "New Channel").firstResult();
        assertNotNull(persistedChannel);
    }

    @Test
    @TestTransaction
    public void testDeleteChannel() {
        // Delete a channel

        proto.DeleteChannelRequest request = proto.DeleteChannelRequest.newBuilder()
                .setName(TEST_CHANNEL_NAME) 
                .build();        
        channelController.deleteChannel(request);

        // Verify the channel was deleted
        Channel deletedChannel = Channel.find("name", TEST_CHANNEL_NAME).firstResult();
        assertNull(deletedChannel);
    }
}
